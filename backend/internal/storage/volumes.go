package storage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// ListVolumes lists all storage volumes/pools
func ListVolumes() ([]Volume, error) {
	var volumes []Volume

	// Get mounted filesystems
	mountedVolumes, err := getMountedVolumes()
	if err != nil {
		logger.Warn("Failed to get mounted volumes", zap.Error(err))
	} else {
		volumes = append(volumes, mountedVolumes...)
	}

	// Get RAID arrays
	raidVolumes, err := getRAIDVolumes()
	if err != nil {
		logger.Warn("Failed to get RAID volumes", zap.Error(err))
	} else {
		volumes = append(volumes, raidVolumes...)
	}

	// Get LVM volumes
	lvmVolumes, err := getLVMVolumes()
	if err != nil {
		logger.Warn("Failed to get LVM volumes", zap.Error(err))
	} else {
		volumes = append(volumes, lvmVolumes...)
	}

	return volumes, nil
}

// getMountedVolumes gets all mounted filesystems
func getMountedVolumes() ([]Volume, error) {
	var volumes []Volume

	cmd := exec.Command("df", "-B1", "--output=source,fstype,size,used,avail,target")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue // Skip header
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		source := fields[0]
		fstype := fields[1]
		size, _ := strconv.ParseUint(fields[2], 10, 64)
		used, _ := strconv.ParseUint(fields[3], 10, 64)
		avail, _ := strconv.ParseUint(fields[4], 10, 64)
		mountPoint := fields[5]

		// Skip special filesystems
		if strings.HasPrefix(source, "tmpfs") ||
			strings.HasPrefix(source, "devtmpfs") ||
			strings.HasPrefix(source, "udev") ||
			strings.HasPrefix(mountPoint, "/sys") ||
			strings.HasPrefix(mountPoint, "/proc") ||
			strings.HasPrefix(mountPoint, "/dev") ||
			strings.HasPrefix(mountPoint, "/run") {
			continue
		}

		volume := Volume{
			ID:         filepath.Base(source),
			Name:       filepath.Base(source),
			Type:       VolumeTypeSingle,
			Status:     VolumeStatusOnline,
			Size:       size,
			Used:       used,
			Available:  avail,
			Filesystem: fstype,
			MountPoint: mountPoint,
			Health:     100,
			CreatedAt:  time.Now(), // We can't easily get the actual creation time
		}

		// Try to get the underlying disk
		if strings.HasPrefix(source, "/dev/") {
			diskName := strings.TrimPrefix(source, "/dev/")
			// Remove partition number
			diskName = strings.TrimRight(diskName, "0123456789")
			if diskName != "" {
				volume.Disks = []string{diskName}
			}
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// getRAIDVolumes gets all RAID arrays
func getRAIDVolumes() ([]Volume, error) {
	var volumes []Volume

	// Check if mdadm is available
	if _, err := exec.LookPath("mdadm"); err != nil {
		return volumes, nil // No RAID support
	}

	// Read /proc/mdstat
	data, err := os.ReadFile("/proc/mdstat")
	if err != nil {
		return volumes, nil // No RAID arrays
	}

	lines := strings.Split(string(data), "\n")
	var currentRaid *Volume

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New RAID array line (e.g., "md0 : active raid1 sda1[0] sdb1[1]")
		if strings.HasPrefix(line, "md") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			raidName := fields[0]
			status := fields[2]
			raidLevel := fields[3]

			// Extract disk members
			var disks []string
			for i := 4; i < len(fields); i++ {
				disk := strings.Split(fields[i], "[")[0]
				disks = append(disks, disk)
			}

			currentRaid = &Volume{
				ID:        raidName,
				Name:      raidName,
				Type:      getRaidVolumeType(raidLevel),
				RaidLevel: raidLevel,
				Disks:     disks,
				CreatedAt: time.Now(),
			}

			if status == "active" {
				currentRaid.Status = VolumeStatusOnline
				currentRaid.Health = 100
			} else {
				currentRaid.Status = VolumeStatusOffline
				currentRaid.Health = 0
			}

			volumes = append(volumes, *currentRaid)
		}

		// Size info line (e.g., "1953383488 blocks super 1.2 [2/2] [UU]")
		if currentRaid != nil && strings.Contains(line, "blocks") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				blocks, _ := strconv.ParseUint(fields[0], 10, 64)
				currentRaid.Size = blocks * 1024 // blocks are in KB

				// Check rebuild status [2/1] means degraded
				if strings.Contains(line, "[U_]") || strings.Contains(line, "[_U]") {
					currentRaid.Status = VolumeStatusDegraded
					currentRaid.Health = 50
				}
			}
		}

		// Rebuild progress
		if currentRaid != nil && strings.Contains(line, "recovery") {
			currentRaid.Status = VolumeStatusRebuilding
		}
	}

	return volumes, nil
}

// getLVMVolumes gets all LVM logical volumes
func getLVMVolumes() ([]Volume, error) {
	var volumes []Volume

	// Check if lvs is available
	if _, err := exec.LookPath("lvs"); err != nil {
		return volumes, nil // No LVM support
	}

	cmd := exec.Command("lvs", "--noheadings", "--units", "B", "-o", "lv_name,vg_name,lv_size,lv_path,lv_attr")
	output, err := cmd.Output()
	if err != nil {
		return volumes, nil // No LVM volumes
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		lvName := fields[0]
		vgName := fields[1]
		sizeStr := strings.TrimSuffix(fields[2], "B")
		lvPath := fields[3]
		attr := fields[4]

		size, _ := strconv.ParseUint(sizeStr, 10, 64)

		// Determine status from attributes
		status := VolumeStatusOnline
		health := 100
		if len(attr) > 4 && attr[4] != 'a' {
			status = VolumeStatusOffline
			health = 0
		}

		volume := Volume{
			ID:        vgName + "/" + lvName,
			Name:      lvName,
			Type:      VolumeTypeLVM,
			Status:    status,
			Size:      size,
			Health:    health,
			CreatedAt: time.Now(),
		}

		// Get filesystem and mount info
		if _, err := os.Stat(lvPath); err == nil {
			cmd := exec.Command("lsblk", "-no", "FSTYPE,MOUNTPOINT", lvPath)
			if output, err := cmd.Output(); err == nil {
				fields := strings.Fields(string(output))
				if len(fields) > 0 {
					volume.Filesystem = fields[0]
				}
				if len(fields) > 1 {
					volume.MountPoint = fields[1]
				}
			}
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// getRaidVolumeType converts RAID level string to VolumeType
func getRaidVolumeType(raidLevel string) VolumeType {
	switch strings.ToLower(raidLevel) {
	case "raid0":
		return VolumeTypeRAID0
	case "raid1":
		return VolumeTypeRAID1
	case "raid5":
		return VolumeTypeRAID5
	case "raid6":
		return VolumeTypeRAID6
	case "raid10":
		return VolumeTypeRAID10
	default:
		return VolumeTypeSingle
	}
}

// GetVolume retrieves information about a specific volume
func GetVolume(id string) (*Volume, error) {
	volumes, err := ListVolumes()
	if err != nil {
		return nil, err
	}

	for _, vol := range volumes {
		if vol.ID == id {
			return &vol, nil
		}
	}

	return nil, fmt.Errorf("volume not found: %s", id)
}

// CreateVolume creates a new storage volume
func CreateVolume(req *CreateVolumeRequest) (*Volume, error) {
	logger.Info("Creating volume",
		zap.String("name", req.Name),
		zap.String("type", string(req.Type)),
		zap.Strings("disks", req.Disks))

	switch req.Type {
	case VolumeTypeRAID0, VolumeTypeRAID1, VolumeTypeRAID5, VolumeTypeRAID6, VolumeTypeRAID10:
		return createRAIDVolume(req)
	case VolumeTypeLVM:
		return createLVMVolume(req)
	case VolumeTypeSingle:
		return createSingleVolume(req)
	default:
		return nil, fmt.Errorf("unsupported volume type: %s", req.Type)
	}
}

// createRAIDVolume creates a RAID array
func createRAIDVolume(req *CreateVolumeRequest) (*Volume, error) {
	if _, err := exec.LookPath("mdadm"); err != nil {
		return nil, fmt.Errorf("mdadm not available")
	}

	// Prepare disk paths
	var diskPaths []string
	for _, disk := range req.Disks {
		if !strings.HasPrefix(disk, "/dev/") {
			disk = "/dev/" + disk
		}
		diskPaths = append(diskPaths, disk)
	}

	// Determine RAID level
	raidLevel := strings.ToLower(string(req.Type))
	raidLevel = strings.TrimPrefix(raidLevel, "raid")

	// Create RAID array
	args := []string{
		"--create",
		"/dev/md/" + req.Name,
		"--level=" + raidLevel,
		"--raid-devices=" + strconv.Itoa(len(diskPaths)),
	}
	args = append(args, diskPaths...)

	cmd := exec.Command("mdadm", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create RAID: %s: %w", string(output), err)
	}

	// Format the RAID array
	formatReq := &FormatDiskRequest{
		Disk:       "/dev/md/" + req.Name,
		Filesystem: req.Filesystem,
		Label:      req.Name,
	}

	if err := FormatDisk(formatReq); err != nil {
		return nil, fmt.Errorf("failed to format RAID: %w", err)
	}

	// Create mount point and mount
	if err := os.MkdirAll(req.MountPoint, 0755); err != nil {
		return nil, fmt.Errorf("failed to create mount point: %w", err)
	}

	cmd = exec.Command("mount", "/dev/md/"+req.Name, req.MountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to mount: %s: %w", string(output), err)
	}

	logger.Info("RAID volume created successfully", zap.String("name", req.Name))

	return GetVolume("md/" + req.Name)
}

// createLVMVolume creates an LVM logical volume
func createLVMVolume(req *CreateVolumeRequest) (*Volume, error) {
	if _, err := exec.LookPath("lvcreate"); err != nil {
		return nil, fmt.Errorf("LVM tools not available")
	}

	// For now, this is a simplified version
	// In production, you'd need to create VG first, then LV
	return nil, fmt.Errorf("LVM volume creation not yet implemented")
}

// createSingleVolume creates a single-disk volume
func createSingleVolume(req *CreateVolumeRequest) (*Volume, error) {
	if len(req.Disks) != 1 {
		return nil, fmt.Errorf("single volume requires exactly one disk")
	}

	disk := req.Disks[0]
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	// Validate disk exists
	if _, err := os.Stat(disk); os.IsNotExist(err) {
		return nil, fmt.Errorf("disk not found: %s", disk)
	}

	// Check if disk already has a filesystem
	existingFS := getExistingFilesystem(disk)
	if existingFS != "" {
		logger.Info("Disk already has a filesystem",
			zap.String("disk", disk),
			zap.String("filesystem", existingFS))

		// If filesystem matches requested, skip formatting
		if existingFS == req.Filesystem {
			logger.Info("Filesystem matches requested type, skipping format",
				zap.String("disk", disk),
				zap.String("filesystem", existingFS))
		} else {
			// Different filesystem, need to format
			logger.Warn("Disk has different filesystem, formatting will destroy data",
				zap.String("disk", disk),
				zap.String("existing", existingFS),
				zap.String("requested", req.Filesystem))

			formatReq := &FormatDiskRequest{
				Disk:       disk,
				Filesystem: req.Filesystem,
				Label:      req.Name,
				Force:      true,
			}

			if err := FormatDisk(formatReq); err != nil {
				return nil, fmt.Errorf("failed to format disk: %w", err)
			}
		}
	} else {
		// No filesystem, format the disk
		logger.Info("No filesystem detected, formatting disk", zap.String("disk", disk))

		formatReq := &FormatDiskRequest{
			Disk:       disk,
			Filesystem: req.Filesystem,
			Label:      req.Name,
		}

		if err := FormatDisk(formatReq); err != nil {
			return nil, fmt.Errorf("failed to format disk: %w", err)
		}
	}

	// Create mount point and mount
	if err := os.MkdirAll(req.MountPoint, 0755); err != nil {
		return nil, fmt.Errorf("failed to create mount point: %w", err)
	}

	// Mount the disk
	cmd := exec.Command("mount", disk, req.MountPoint)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to mount: %s: %w", string(output), err)
	}

	// Add to /etc/fstab for persistence (optional, could be dangerous in dev)
	logger.Info("Skipping /etc/fstab update - not safe for development")
	// addToFstab(disk, req.MountPoint, req.Filesystem)

	logger.Info("Single volume created successfully", zap.String("name", req.Name), zap.String("mount", req.MountPoint))

	return GetVolume(filepath.Base(disk))
}

// getExistingFilesystem checks if a disk/partition has an existing filesystem
func getExistingFilesystem(disk string) string {
	cmd := exec.Command("blkid", "-s", "TYPE", "-o", "value", disk)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// DeleteVolume deletes a storage volume
func DeleteVolume(id string) error {
	volume, err := GetVolume(id)
	if err != nil {
		return err
	}

	// Unmount if mounted
	if volume.MountPoint != "" {
		cmd := exec.Command("umount", volume.MountPoint)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to unmount: %s: %w", string(output), err)
		}
	}

	// Delete based on type
	switch volume.Type {
	case VolumeTypeRAID0, VolumeTypeRAID1, VolumeTypeRAID5, VolumeTypeRAID6, VolumeTypeRAID10:
		cmd := exec.Command("mdadm", "--stop", "/dev/"+id)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to stop RAID: %s: %w", string(output), err)
		}

		// Zero superblocks
		for _, disk := range volume.Disks {
			exec.Command("mdadm", "--zero-superblock", "/dev/"+disk).Run()
		}

	case VolumeTypeLVM:
		cmd := exec.Command("lvremove", "-f", id)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to remove LV: %s: %w", string(output), err)
		}
	}

	// Remove from fstab
	removeFromFstab(volume.MountPoint)

	logger.Info("Volume deleted successfully", zap.String("id", id))

	return nil
}

// addToFstab adds an entry to /etc/fstab
func addToFstab(device, mountPoint, fstype string) error {
	// Get UUID
	cmd := exec.Command("blkid", "-s", "UUID", "-o", "value", device)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	uuid := strings.TrimSpace(string(output))
	if uuid == "" {
		return fmt.Errorf("failed to get UUID")
	}

	entry := fmt.Sprintf("UUID=%s %s %s defaults 0 2\n", uuid, mountPoint, fstype)

	file, err := os.OpenFile("/etc/fstab", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(entry)
	return err
}

// removeFromFstab removes an entry from /etc/fstab
func removeFromFstab(mountPoint string) error {
	// This is a simplified version
	// In production, you'd want to parse and rewrite /etc/fstab properly
	return nil
}
