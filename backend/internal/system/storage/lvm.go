// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package storage

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"fmt"
	"strconv"
	"strings"
)

// LVMManager manages LVM (Logical Volume Manager)
type LVMManager struct {
	shell      executor.ShellExecutor
	enabled bool
}

// VolumeGroup represents an LVM volume group
type VolumeGroup struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	Size       uint64 `json:"size"`
	Free       uint64 `json:"free"`
	PVCount    int    `json:"pv_count"`
	LVCount    int    `json:"lv_count"`
	SnapshotCount int `json:"snapshot_count"`
}

// LogicalVolume represents an LVM logical volume
type LogicalVolume struct {
	Name       string `json:"name"`
	VGName     string `json:"vg_name"`
	UUID       string `json:"uuid"`
	Size       uint64 `json:"size"`
	Path       string `json:"path"`
	Active     bool   `json:"active"`
}

// NewLVMManager creates a new LVM manager
func NewLVMManager(shell executor.ShellExecutor) (*LVMManager, error) {
	if !shell.CommandExists("vgs") || !shell.CommandExists("lvs") {
		return nil, fmt.Errorf("LVM tools not installed")
	}

	return &LVMManager{
		shell:   shell,
		enabled: true,
	}, nil
}

// IsEnabled returns whether LVM is available
func (l *LVMManager) IsEnabled() bool {
	return l.enabled
}

// ListVolumeGroups lists all volume groups
func (l *LVMManager) ListVolumeGroups() ([]VolumeGroup, error) {
	result, err := l.shell.Execute("vgs", "--noheadings", "--units", "b", "--separator", "|",
		"-o", "vg_name,vg_uuid,vg_size,vg_free,pv_count,lv_count,snap_count")
	if err != nil {
		return nil, fmt.Errorf("failed to list volume groups: %w", err)
	}

	var vgs []VolumeGroup
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(strings.TrimSpace(line), "|")
		if len(fields) < 7 {
			continue
		}

		vg := VolumeGroup{
			Name: strings.TrimSpace(fields[0]),
			UUID: strings.TrimSpace(fields[1]),
		}

		// Parse sizes (remove 'B' suffix)
		if size, err := parseSize(fields[2]); err == nil {
			vg.Size = size
		}
		if free, err := parseSize(fields[3]); err == nil {
			vg.Free = free
		}

		// Parse counts
		if count, err := strconv.Atoi(strings.TrimSpace(fields[4])); err == nil {
			vg.PVCount = count
		}
		if count, err := strconv.Atoi(strings.TrimSpace(fields[5])); err == nil {
			vg.LVCount = count
		}
		if count, err := strconv.Atoi(strings.TrimSpace(fields[6])); err == nil {
			vg.SnapshotCount = count
		}

		vgs = append(vgs, vg)
	}

	return vgs, nil
}

// ListLogicalVolumes lists all logical volumes
func (l *LVMManager) ListLogicalVolumes(vgName string) ([]LogicalVolume, error) {
	args := []string{"--noheadings", "--units", "b", "--separator", "|",
		"-o", "lv_name,vg_name,lv_uuid,lv_size,lv_path,lv_active"}

	if vgName != "" {
		args = append(args, vgName)
	}

	result, err := l.shell.Execute("lvs", args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list logical volumes: %w", err)
	}

	var lvs []LogicalVolume
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(strings.TrimSpace(line), "|")
		if len(fields) < 6 {
			continue
		}

		lv := LogicalVolume{
			Name:   strings.TrimSpace(fields[0]),
			VGName: strings.TrimSpace(fields[1]),
			UUID:   strings.TrimSpace(fields[2]),
			Path:   strings.TrimSpace(fields[4]),
			Active: strings.TrimSpace(fields[5]) == "active",
		}

		if size, err := parseSize(fields[3]); err == nil {
			lv.Size = size
		}

		lvs = append(lvs, lv)
	}

	return lvs, nil
}

// CreateVolumeGroup creates a new volume group
func (l *LVMManager) CreateVolumeGroup(name string, devices []string) error {
	if len(devices) == 0 {
		return fmt.Errorf("no devices specified")
	}

	args := append([]string{name}, devices...)
	_, err := l.shell.Execute("vgcreate", args...)
	if err != nil {
		return fmt.Errorf("failed to create volume group: %w", err)
	}

	return nil
}

// CreateLogicalVolume creates a new logical volume
func (l *LVMManager) CreateLogicalVolume(vgName string, lvName string, size string) error {
	_, err := l.shell.Execute("lvcreate", "-L", size, "-n", lvName, vgName)
	if err != nil {
		return fmt.Errorf("failed to create logical volume: %w", err)
	}

	return nil
}

// DeleteLogicalVolume deletes a logical volume
func (l *LVMManager) DeleteLogicalVolume(vgName string, lvName string, force bool) error {
	args := []string{fmt.Sprintf("%s/%s", vgName, lvName)}
	if force {
		args = append([]string{"-f"}, args...)
	}

	_, err := l.shell.Execute("lvremove", args...)
	if err != nil {
		return fmt.Errorf("failed to delete logical volume: %w", err)
	}

	return nil
}

// ExtendLogicalVolume extends a logical volume
func (l *LVMManager) ExtendLogicalVolume(vgName string, lvName string, size string) error {
	_, err := l.shell.Execute("lvextend", "-L", "+"+size, fmt.Sprintf("%s/%s", vgName, lvName))
	if err != nil {
		return fmt.Errorf("failed to extend logical volume: %w", err)
	}

	return nil
}

// CreateSnapshot creates an LVM snapshot
func (l *LVMManager) CreateSnapshot(vgName string, lvName string, snapshotName string, size string) error {
	_, err := l.shell.Execute("lvcreate", "-L", size, "-s", "-n", snapshotName,
		fmt.Sprintf("%s/%s", vgName, lvName))
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	return nil
}

// Helper to parse size strings like "10.00G" or "1024B"
func parseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	sizeStr = strings.TrimSuffix(sizeStr, "B")

	return strconv.ParseUint(sizeStr, 10, 64)
}
