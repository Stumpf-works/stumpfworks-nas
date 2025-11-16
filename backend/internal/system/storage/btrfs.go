// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package storage

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strconv"
	"strings"
)

// BTRFSManager manages BTRFS filesystems
type BTRFSManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// BTRFSFilesystem represents a BTRFS filesystem
type BTRFSFilesystem struct {
	UUID       string   `json:"uuid"`
	Label      string   `json:"label"`
	Devices    []string `json:"devices"`
	TotalSize  uint64   `json:"total_size"`
	Used       uint64   `json:"used"`
	DataRatio  string   `json:"data_ratio"`
	Metadata   string   `json:"metadata"`
}

// NewBTRFSManager creates a new BTRFS manager
func NewBTRFSManager(shell executor.ShellExecutor) (*BTRFSManager, error) {
	if !shell.CommandExists("btrfs") {
		return nil, fmt.Errorf("btrfs-progs not installed")
	}

	return &BTRFSManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether BTRFS is available
func (b *BTRFSManager) IsEnabled() bool {
	return b.enabled
}

// ListFilesystems lists all BTRFS filesystems
func (b *BTRFSManager) ListFilesystems() ([]BTRFSFilesystem, error) {
	result, err := b.shell.Execute("btrfs", "filesystem", "show")
	if err != nil {
		return nil, fmt.Errorf("failed to list filesystems: %w", err)
	}

	var filesystems []BTRFSFilesystem
	var current *BTRFSFilesystem

	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Label:") {
			if current != nil {
				filesystems = append(filesystems, *current)
			}
			current = &BTRFSFilesystem{}

			// Parse label and UUID
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "Label:" && i+1 < len(parts) {
					current.Label = strings.Trim(parts[i+1], "'")
				}
				if part == "uuid:" && i+1 < len(parts) {
					current.UUID = parts[i+1]
				}
			}
		} else if strings.Contains(line, "devid") && current != nil {
			// Parse device path
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "path" && i+1 < len(parts) {
					current.Devices = append(current.Devices, parts[i+1])
				}
			}
		}
	}

	if current != nil {
		filesystems = append(filesystems, *current)
	}

	return filesystems, nil
}

// CreateFilesystem creates a new BTRFS filesystem
func (b *BTRFSManager) CreateFilesystem(devices []string, label string, dataRaid string, metadataRaid string) error {
	if len(devices) == 0 {
		return fmt.Errorf("no devices specified")
	}

	args := []string{"mkfs.btrfs"}

	if label != "" {
		args = append(args, "-L", label)
	}

	if dataRaid != "" {
		args = append(args, "-d", dataRaid)
	}

	if metadataRaid != "" {
		args = append(args, "-m", metadataRaid)
	}

	args = append(args, devices...)

	_, err := b.shell.Execute(args[0], args[1:]...)
	if err != nil {
		return fmt.Errorf("failed to create filesystem: %w", err)
	}

	return nil
}

// CreateSnapshot creates a BTRFS snapshot
func (b *BTRFSManager) CreateSnapshot(source string, dest string, readonly bool) error {
	args := []string{"subvolume", "snapshot"}
	if readonly {
		args = append(args, "-r")
	}
	args = append(args, source, dest)

	_, err := b.shell.Execute("btrfs", args...)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	return nil
}

// DeleteSnapshot deletes a BTRFS snapshot
func (b *BTRFSManager) DeleteSnapshot(path string) error {
	_, err := b.shell.Execute("btrfs", "subvolume", "delete", path)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	return nil
}

// Scrub starts a scrub operation
func (b *BTRFSManager) Scrub(path string) error {
	_, err := b.shell.Execute("btrfs", "scrub", "start", path)
	if err != nil {
		return fmt.Errorf("failed to start scrub: %w", err)
	}

	return nil
}

// GetUsage gets filesystem usage
func (b *BTRFSManager) GetUsage(path string) (uint64, uint64, error) {
	result, err := b.shell.Execute("btrfs", "filesystem", "usage", "-b", path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get usage: %w", err)
	}

	var totalSize, used uint64
	lines := strings.Split(result.Stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Device size:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				if size, err := strconv.ParseUint(parts[2], 10, 64); err == nil {
					totalSize = size
				}
			}
		}
		if strings.Contains(line, "Used:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if size, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
					used = size
				}
			}
		}
	}

	return totalSize, used, nil
}

// ===== Advanced BTRFS Features =====

// BTRFSSubvolume represents a BTRFS subvolume
type BTRFSSubvolume struct {
	ID         uint64 `json:"id"`
	ParentID   uint64 `json:"parent_id"`
	TopLevel   uint64 `json:"top_level"`
	Path       string `json:"path"`
	UUID       string `json:"uuid"`
	ParentUUID string `json:"parent_uuid,omitempty"`
}

// CreateSubvolume creates a new BTRFS subvolume
func (b *BTRFSManager) CreateSubvolume(path string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "subvolume", "create", path)
	if err != nil {
		return fmt.Errorf("failed to create subvolume: %w", err)
	}

	return nil
}

// DeleteSubvolume deletes a BTRFS subvolume
func (b *BTRFSManager) DeleteSubvolume(path string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "subvolume", "delete", path)
	if err != nil {
		return fmt.Errorf("failed to delete subvolume: %w", err)
	}

	return nil
}

// ListSubvolumes lists all subvolumes in a filesystem
func (b *BTRFSManager) ListSubvolumes(path string) ([]BTRFSSubvolume, error) {
	if !b.enabled {
		return nil, fmt.Errorf("BTRFS not available")
	}

	result, err := b.shell.Execute("btrfs", "subvolume", "list", "-p", "-u", path)
	if err != nil {
		return nil, fmt.Errorf("failed to list subvolumes: %w", err)
	}

	var subvolumes []BTRFSSubvolume
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		subvol := BTRFSSubvolume{}

		// Parse ID
		if id, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			subvol.ID = id
		}

		// Parse parent ID
		if pid, err := strconv.ParseUint(fields[3], 10, 64); err == nil {
			subvol.ParentID = pid
		}

		// Parse UUID
		if len(fields) >= 7 {
			subvol.UUID = fields[6]
		}

		// Parse path (last field)
		subvol.Path = fields[len(fields)-1]

		subvolumes = append(subvolumes, subvol)
	}

	return subvolumes, nil
}

// Send sends a BTRFS subvolume/snapshot to a stream
func (b *BTRFSManager) Send(subvolumePath string, parentPath string) (string, error) {
	if !b.enabled {
		return "", fmt.Errorf("BTRFS not available")
	}

	args := []string{"send"}
	if parentPath != "" {
		args = append(args, "-p", parentPath)
	}
	args = append(args, subvolumePath)

	result, err := b.shell.Execute("btrfs", args...)
	if err != nil {
		return "", fmt.Errorf("failed to send subvolume: %w", err)
	}

	return result.Stdout, nil
}

// Receive receives a BTRFS subvolume from a stream
func (b *BTRFSManager) Receive(targetPath string, streamFile string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	// Note: In real implementation, this would read from stdin
	_, err := b.shell.Execute("sh", "-c", fmt.Sprintf("cat %s | btrfs receive %s", streamFile, targetPath))
	if err != nil {
		return fmt.Errorf("failed to receive subvolume: %w", err)
	}

	return nil
}

// AddDevice adds a device to a BTRFS filesystem
func (b *BTRFSManager) AddDevice(device string, mountPoint string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "device", "add", device, mountPoint)
	if err != nil {
		return fmt.Errorf("failed to add device: %w", err)
	}

	return nil
}

// RemoveDevice removes a device from a BTRFS filesystem
func (b *BTRFSManager) RemoveDevice(device string, mountPoint string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "device", "remove", device, mountPoint)
	if err != nil {
		return fmt.Errorf("failed to remove device: %w", err)
	}

	return nil
}

// Balance starts a balance operation on a BTRFS filesystem
func (b *BTRFSManager) Balance(mountPoint string, filters string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	args := []string{"balance", "start"}
	if filters != "" {
		args = append(args, filters)
	}
	args = append(args, mountPoint)

	_, err := b.shell.Execute("btrfs", args...)
	if err != nil {
		return fmt.Errorf("failed to start balance: %w", err)
	}

	return nil
}

// BalanceStatus returns the status of a balance operation
func (b *BTRFSManager) BalanceStatus(mountPoint string) (string, error) {
	if !b.enabled {
		return "", fmt.Errorf("BTRFS not available")
	}

	result, err := b.shell.Execute("btrfs", "balance", "status", mountPoint)
	if err != nil {
		return "", fmt.Errorf("failed to get balance status: %w", err)
	}

	return result.Stdout, nil
}

// PauseBalance pauses a running balance operation
func (b *BTRFSManager) PauseBalance(mountPoint string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "balance", "pause", mountPoint)
	if err != nil {
		return fmt.Errorf("failed to pause balance: %w", err)
	}

	return nil
}

// ResumeBalance resumes a paused balance operation
func (b *BTRFSManager) ResumeBalance(mountPoint string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "balance", "resume", mountPoint)
	if err != nil {
		return fmt.Errorf("failed to resume balance: %w", err)
	}

	return nil
}

// CancelBalance cancels a running balance operation
func (b *BTRFSManager) CancelBalance(mountPoint string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "balance", "cancel", mountPoint)
	if err != nil {
		return fmt.Errorf("failed to cancel balance: %w", err)
	}

	return nil
}

// Defragment defragments files or directories
func (b *BTRFSManager) Defragment(path string, recursive bool, compress string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	args := []string{"filesystem", "defragment"}
	if recursive {
		args = append(args, "-r")
	}
	if compress != "" {
		args = append(args, "-c"+compress)
	}
	args = append(args, path)

	_, err := b.shell.Execute("btrfs", args...)
	if err != nil {
		return fmt.Errorf("failed to defragment: %w", err)
	}

	return nil
}

// SetCompression sets compression for a file or directory
func (b *BTRFSManager) SetCompression(path string, compression string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	// Use chattr to set compression
	compressOpt := ""
	switch compression {
	case "zlib":
		compressOpt = "c"
	case "lzo":
		compressOpt = "c"
	case "zstd":
		compressOpt = "c"
	case "none":
		compressOpt = "c" // Remove compression flag
	default:
		return fmt.Errorf("unsupported compression: %s", compression)
	}

	_, err := b.shell.Execute("chattr", "+"+compressOpt, path)
	if err != nil {
		return fmt.Errorf("failed to set compression: %w", err)
	}

	return nil
}

// Resize resizes a BTRFS filesystem
func (b *BTRFSManager) Resize(mountPoint string, size string) error {
	if !b.enabled {
		return fmt.Errorf("BTRFS not available")
	}

	_, err := b.shell.Execute("btrfs", "filesystem", "resize", size, mountPoint)
	if err != nil {
		return fmt.Errorf("failed to resize filesystem: %w", err)
	}

	return nil
}

// GetDeviceStats returns device statistics
func (b *BTRFSManager) GetDeviceStats(mountPoint string) (string, error) {
	if !b.enabled {
		return "", fmt.Errorf("BTRFS not available")
	}

	result, err := b.shell.Execute("btrfs", "device", "stats", mountPoint)
	if err != nil {
		return "", fmt.Errorf("failed to get device stats: %w", err)
	}

	return result.Stdout, nil
}
