package storage

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// ListDisks lists all available disks on the system
func ListDisks() ([]Disk, error) {
	var disks []Disk

	// Read from /sys/block to get all block devices
	blockDevices, err := filepath.Glob("/sys/block/*")
	if err != nil {
		return nil, fmt.Errorf("failed to read block devices: %w", err)
	}

	for _, blockPath := range blockDevices {
		diskName := filepath.Base(blockPath)

		// Skip loop devices, ram devices, etc.
		if strings.HasPrefix(diskName, "loop") ||
			strings.HasPrefix(diskName, "ram") ||
			strings.HasPrefix(diskName, "dm-") {
			continue
		}

		disk, err := GetDiskInfo(diskName)
		if err != nil {
			logger.Warn("Failed to get disk info", zap.String("disk", diskName), zap.Error(err))
			continue
		}

		disks = append(disks, *disk)
	}

	return disks, nil
}

// GetDiskInfo retrieves detailed information about a specific disk
func GetDiskInfo(diskName string) (*Disk, error) {
	diskPath := "/dev/" + diskName
	sysPath := "/sys/block/" + diskName

	// Check if disk exists
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("disk not found: %s", diskName)
	}

	disk := &Disk{
		Name: diskName,
		Path: diskPath,
	}

	// Get disk size
	size, err := getDiskSize(sysPath)
	if err != nil {
		logger.Warn("Failed to get disk size", zap.String("disk", diskName), zap.Error(err))
	} else {
		disk.Size = size
	}

	// Get disk model
	model, err := readSysFile(filepath.Join(sysPath, "device/model"))
	if err == nil {
		disk.Model = strings.TrimSpace(model)
	}

	// Get disk type (HDD/SSD/NVMe)
	disk.Type = getDiskType(diskName, sysPath)

	// Check if removable
	removable, err := readSysFile(filepath.Join(sysPath, "removable"))
	if err == nil {
		disk.IsRemovable = strings.TrimSpace(removable) == "1"
	}

	// Check if system disk (contains root partition)
	disk.IsSystem = isSystemDisk(diskName)

	// Get partitions
	partitions, err := getPartitions(diskName)
	if err != nil {
		logger.Warn("Failed to get partitions", zap.String("disk", diskName), zap.Error(err))
	} else {
		disk.Partitions = partitions
	}

	// Get SMART data if available
	smart, err := GetSMARTData(diskName)
	if err == nil {
		disk.SMART = smart
		disk.SMARTEnabled = true
		disk.Status = getHealthStatus(smart)
		disk.Temperature = smart.Temperature
	} else {
		disk.Status = DiskStatusUnknown
	}

	return disk, nil
}

// getDiskSize reads the disk size from sysfs
func getDiskSize(sysPath string) (uint64, error) {
	sizeStr, err := readSysFile(filepath.Join(sysPath, "size"))
	if err != nil {
		return 0, err
	}

	// Size is in 512-byte sectors
	sectors, err := strconv.ParseUint(strings.TrimSpace(sizeStr), 10, 64)
	if err != nil {
		return 0, err
	}

	return sectors * 512, nil
}

// getDiskType determines the type of disk
func getDiskType(diskName, sysPath string) DiskType {
	if strings.HasPrefix(diskName, "nvme") {
		return DiskTypeNVMe
	}

	// Check if it's an SSD by checking rotational
	rotational, err := readSysFile(filepath.Join(sysPath, "queue/rotational"))
	if err == nil && strings.TrimSpace(rotational) == "0" {
		return DiskTypeSSD
	}

	// Check if USB
	if strings.Contains(sysPath, "usb") {
		return DiskTypeUSB
	}

	return DiskTypeHDD
}

// isSystemDisk checks if the disk contains the root partition
func isSystemDisk(diskName string) bool {
	// Read /proc/mounts
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		device := fields[0]
		mountPoint := fields[1]

		// Check if this is the root partition
		if mountPoint == "/" && strings.Contains(device, diskName) {
			return true
		}
	}

	return false
}

// getPartitions retrieves all partitions for a disk
func getPartitions(diskName string) ([]Partition, error) {
	var partitions []Partition

	// Use lsblk to get all partitions for this disk
	// This works for all disk types (sda, nvme0n1, mmcblk0, etc.)
	cmd := exec.Command("lsblk", "-ln", "-o", "NAME,TYPE", "/dev/"+diskName)
	output, err := cmd.Output()
	if err != nil {
		// Fallback to glob pattern for older systems without lsblk
		logger.Warn("lsblk failed, falling back to glob pattern", zap.String("disk", diskName), zap.Error(err))
		return getPartitionsViaGlob(diskName)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		partName := fields[0]
		partType := fields[1]

		// Only include partitions, skip the disk itself
		if partType != "part" {
			continue
		}

		partition, err := getPartitionInfo(partName)
		if err != nil {
			logger.Warn("Failed to get partition info", zap.String("partition", partName), zap.Error(err))
			continue
		}

		partitions = append(partitions, *partition)
	}

	return partitions, nil
}

// getPartitionsViaGlob is a fallback method using glob patterns
func getPartitionsViaGlob(diskName string) ([]Partition, error) {
	var partitions []Partition

	sysPath := filepath.Join("/sys/block", diskName)

	// Try different partition naming patterns
	patterns := []string{
		diskName + "[0-9]*",     // sda1, sda2, sdb1, etc.
		diskName + "p[0-9]*",    // nvme0n1p1, mmcblk0p1, etc.
	}

	for _, pattern := range patterns {
		partPaths, err := filepath.Glob(filepath.Join(sysPath, pattern))
		if err != nil {
			continue
		}

		for _, partPath := range partPaths {
			partName := filepath.Base(partPath)
			if partName == diskName {
				continue // Skip the disk itself
			}

			partition, err := getPartitionInfo(partName)
			if err != nil {
				logger.Warn("Failed to get partition info", zap.String("partition", partName), zap.Error(err))
				continue
			}

			partitions = append(partitions, *partition)
		}
	}

	return partitions, nil
}

// getPartitionInfo retrieves information about a specific partition
func getPartitionInfo(partName string) (*Partition, error) {
	partPath := "/dev/" + partName

	partition := &Partition{
		Name: partName,
		Path: partPath,
	}

	// Get size
	sysPath := findPartitionSysPath(partName)
	if sysPath != "" {
		size, err := getDiskSize(sysPath)
		if err == nil {
			partition.Size = size
		}
	}

	// Use lsblk to get filesystem, mount point, UUID, and label
	cmd := exec.Command("lsblk", "-no", "FSTYPE,MOUNTPOINT,LABEL,UUID", partPath)
	output, err := cmd.Output()
	if err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			partition.Filesystem = fields[0]
		}
		if len(fields) > 1 && fields[1] != "" {
			partition.MountPoint = fields[1]
			partition.IsMounted = true
		}
		if len(fields) > 2 {
			partition.Label = fields[2]
		}
		if len(fields) > 3 {
			partition.UUID = fields[3]
		}
	}

	// Get used space if mounted
	if partition.IsMounted && partition.MountPoint != "" {
		used, err := getUsedSpace(partition.MountPoint)
		if err == nil {
			partition.Used = used
		}
	}

	return partition, nil
}

// findPartitionSysPath finds the sysfs path for a partition
func findPartitionSysPath(partName string) string {
	// Try common patterns
	patterns := []string{
		"/sys/block/*/",
		"/sys/block/*/*/",
	}

	for _, pattern := range patterns {
		path := filepath.Join(pattern, partName)
		matches, _ := filepath.Glob(path)
		if len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// getUsedSpace gets the used space for a mounted filesystem
func getUsedSpace(mountPoint string) (uint64, error) {
	cmd := exec.Command("df", "-B1", "--output=used", mountPoint)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unexpected df output")
	}

	used, err := strconv.ParseUint(strings.TrimSpace(lines[1]), 10, 64)
	if err != nil {
		return 0, err
	}

	return used, nil
}

// GetSMARTData retrieves SMART monitoring data for a disk
func GetSMARTData(diskName string) (*SMARTData, error) {
	diskPath := "/dev/" + diskName

	// Check if smartctl is available
	if _, err := exec.LookPath("smartctl"); err != nil {
		return nil, fmt.Errorf("smartctl not available")
	}

	// Run smartctl
	cmd := exec.Command("smartctl", "-A", "-H", "-i", diskPath)
	output, err := cmd.Output()
	if err != nil {
		// smartctl returns non-zero exit code even for warnings
		// Continue parsing if we got output
		if len(output) == 0 {
			return nil, fmt.Errorf("failed to get SMART data: %w", err)
		}
	}

	smart := &SMARTData{
		LastUpdated: time.Now(),
	}

	outputStr := string(output)

	// Parse health status
	if strings.Contains(outputStr, "PASSED") {
		smart.Healthy = true
	}

	// Parse temperature
	tempRe := regexp.MustCompile(`Temperature_Celsius.*\s+(\d+)`)
	if matches := tempRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if temp, err := strconv.Atoi(matches[1]); err == nil {
			smart.Temperature = temp
		}
	}

	// Parse power on hours
	pohRe := regexp.MustCompile(`Power_On_Hours.*\s+(\d+)`)
	if matches := pohRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if poh, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			smart.PowerOnHours = poh
		}
	}

	// Parse power cycle count
	pccRe := regexp.MustCompile(`Power_Cycle_Count.*\s+(\d+)`)
	if matches := pccRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if pcc, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			smart.PowerCycleCount = pcc
		}
	}

	// Parse reallocated sectors
	rsRe := regexp.MustCompile(`Reallocated_Sector_Ct.*\s+(\d+)`)
	if matches := rsRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if rs, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			smart.ReallocatedSectors = rs
		}
	}

	// Parse pending sectors
	psRe := regexp.MustCompile(`Current_Pending_Sector.*\s+(\d+)`)
	if matches := psRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if ps, err := strconv.ParseUint(matches[1], 10, 64); err == nil {
			smart.PendingSectors = ps
		}
	}

	// Parse SSD wear level
	wearRe := regexp.MustCompile(`Wear_Leveling_Count.*\s+(\d+)`)
	if matches := wearRe.FindStringSubmatch(outputStr); len(matches) > 1 {
		if wear, err := strconv.Atoi(matches[1]); err == nil {
			smart.PercentLifeUsed = 100 - wear
		}
	}

	return smart, nil
}

// getHealthStatus determines the health status based on SMART data
func getHealthStatus(smart *SMARTData) DiskStatus {
	if !smart.Healthy {
		return DiskStatusFailed
	}

	// Critical conditions
	if smart.ReallocatedSectors > 10 ||
		smart.PendingSectors > 5 ||
		smart.UncorrectableErrors > 0 {
		return DiskStatusCritical
	}

	// Warning conditions
	if smart.ReallocatedSectors > 0 ||
		smart.PendingSectors > 0 ||
		smart.Temperature > 60 ||
		smart.PercentLifeUsed > 90 {
		return DiskStatusWarning
	}

	return DiskStatusHealthy
}

// readSysFile reads a single-line file from sysfs
func readSysFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// findSystemBinary searches for a binary in common system paths
func findSystemBinary(name string) (string, error) {
	// First try exec.LookPath (searches in PATH)
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	// Common system paths where binaries like mkfs.* are located
	systemPaths := []string{
		"/usr/sbin",
		"/sbin",
		"/usr/bin",
		"/bin",
		"/usr/local/sbin",
		"/usr/local/bin",
	}

	for _, dir := range systemPaths {
		fullPath := filepath.Join(dir, name)
		if _, err := os.Stat(fullPath); err == nil {
			// Check if executable
			if info, err := os.Stat(fullPath); err == nil {
				if info.Mode()&0111 != 0 {
					return fullPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("%s not found in system paths", name)
}

// GetStorageStats returns overall storage statistics
func GetStorageStats() (*StorageStats, error) {
	disks, err := ListDisks()
	if err != nil {
		return nil, err
	}

	stats := &StorageStats{}

	for _, disk := range disks {
		stats.TotalDisks++
		stats.TotalCapacity += disk.Size

		for _, part := range disk.Partitions {
			stats.UsedCapacity += part.Used
		}

		switch disk.Status {
		case DiskStatusHealthy:
			stats.HealthyDisks++
		case DiskStatusWarning:
			stats.WarningDisks++
		case DiskStatusCritical, DiskStatusFailed:
			stats.CriticalDisks++
		}
	}

	stats.AvailableCapacity = stats.TotalCapacity - stats.UsedCapacity

	return stats, nil
}

// FormatDisk formats a disk or partition with the specified filesystem
func FormatDisk(req *FormatDiskRequest) error {
	diskPath := req.Disk
	if !strings.HasPrefix(diskPath, "/dev/") {
		diskPath = "/dev/" + diskPath
	}

	// Check if disk exists
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		return fmt.Errorf("disk not found: %s", diskPath)
	}

	// Check if mounted (unless force)
	if !req.Force {
		cmd := exec.Command("findmnt", "-n", "-S", diskPath)
		if err := cmd.Run(); err == nil {
			return fmt.Errorf("disk is mounted, use force to format anyway")
		}
	}

	// Unmount if mounted and force is true
	if req.Force {
		exec.Command("umount", diskPath).Run()
	}

	// Determine the mkfs command based on filesystem type
	var mkfsBinary string
	var args []string

	switch req.Filesystem {
	case "ext4":
		mkfsBinary = "mkfs.ext4"
		args = []string{"-F"}
		if req.Label != "" {
			args = append(args, "-L", req.Label)
		}
	case "xfs":
		mkfsBinary = "mkfs.xfs"
		args = []string{"-f"}
		if req.Label != "" {
			args = append(args, "-L", req.Label)
		}
	case "btrfs":
		mkfsBinary = "mkfs.btrfs"
		args = []string{"-f"}
		if req.Label != "" {
			args = append(args, "-L", req.Label)
		}
	default:
		return fmt.Errorf("unsupported filesystem: %s", req.Filesystem)
	}

	// Find the full path to the mkfs tool (checking common system paths)
	mkfsCmd, err := findSystemBinary(mkfsBinary)
	if err != nil {
		logger.Warn("Filesystem tool not available",
			zap.String("tool", mkfsBinary),
			zap.String("disk", diskPath),
			zap.String("filesystem", req.Filesystem))
		return fmt.Errorf("filesystem tool not available: %s is not installed on this system. Please install the required packages (e.g., e2fsprogs for ext4, xfsprogs for xfs, btrfs-progs for btrfs)", mkfsBinary)
	}

	logger.Info("Using filesystem tool", zap.String("path", mkfsCmd))

	// Format the disk
	args = append(args, diskPath)
	cmd := exec.Command(mkfsCmd, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("format failed: %s: %w", string(output), err)
	}

	logger.Info("Disk formatted successfully",
		zap.String("disk", diskPath),
		zap.String("filesystem", req.Filesystem))

	return nil
}
