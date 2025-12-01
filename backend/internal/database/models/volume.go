// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package models

import (
	"time"
)

// Volume represents a storage volume/pool in the database
// This ensures volumes persist across reboots and can be automatically mounted on startup
type Volume struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Type        string    `gorm:"not null" json:"type"` // single, raid0, raid1, raid5, raid6, raid10, lvm
	MountPoint  string    `gorm:"uniqueIndex;not null" json:"mount_point"`
	Disks       string    `gorm:"type:text" json:"disks"` // Comma-separated list of disk names (e.g., "sda,sdb,sdc")
	Filesystem  string    `gorm:"not null" json:"filesystem"` // ext4, xfs, btrfs, etc.
	Status      string    `gorm:"default:online" json:"status"` // online, offline, degraded, error, mounting
	Size        uint64    `json:"size"`        // Total size in bytes (updated periodically)
	Used        uint64    `json:"used"`        // Used space in bytes (updated periodically)
	Available   uint64    `json:"available"`   // Available space in bytes (updated periodically)
	LastMounted *time.Time `json:"last_mounted,omitempty"` // Last successful mount time
	LastError   string    `gorm:"type:text" json:"last_error,omitempty"` // Last mount error message
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Volume model
func (Volume) TableName() string {
	return "volumes"
}
