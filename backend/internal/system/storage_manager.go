// Revision: 2025-11-16 | Author: StumpfWorks AI | Version: 1.1.0
package system

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/storage"
)

// StorageManager manages all storage-related operations
type StorageManager struct {
	shell *ShellExecutor

	// Subsystems
	ZFS   *storage.ZFSManager
	BTRFS *storage.BTRFSManager
	LVM   *storage.LVMManager
	RAID  *storage.RAIDManager
	SMART *storage.SMARTManager
}

// NewStorageManager creates a new storage manager
func NewStorageManager(shell *ShellExecutor) (*StorageManager, error) {
	sm := &StorageManager{
		shell: shell,
	}

	// Initialize ZFS manager
	zfs, err := storage.NewZFSManager(shell)
	if err != nil {
		// ZFS is optional
		sm.ZFS = nil
	} else {
		sm.ZFS = zfs
	}

	// Initialize BTRFS manager
	btrfs, err := storage.NewBTRFSManager(shell)
	if err != nil {
		// BTRFS is optional
		sm.BTRFS = nil
	} else {
		sm.BTRFS = btrfs
	}

	// Initialize LVM manager
	lvm, err := storage.NewLVMManager(shell)
	if err != nil {
		// LVM is optional
		sm.LVM = nil
	} else {
		sm.LVM = lvm
	}

	// Initialize RAID manager
	raid, err := storage.NewRAIDManager(shell)
	if err != nil {
		// RAID is optional
		sm.RAID = nil
	} else {
		sm.RAID = raid
	}

	// Initialize SMART manager
	smart, err := storage.NewSMARTManager(shell)
	if err != nil {
		return nil, fmt.Errorf("SMART manager is required: %w", err)
	}
	sm.SMART = smart

	return sm, nil
}
