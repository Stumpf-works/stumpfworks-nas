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
