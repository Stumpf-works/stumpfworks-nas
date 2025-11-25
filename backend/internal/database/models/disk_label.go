package models

import (
	"gorm.io/gorm"
)

// DiskLabel stores user-defined labels for physical disks
// Key is the disk serial number which is unique and persistent across reboots
type DiskLabel struct {
	gorm.Model
	Serial string `gorm:"size:100;not null;uniqueIndex"` // Disk serial number (unique identifier)
	Label  string `gorm:"size:255;not null"`             // User-defined friendly name
}

// TableName specifies the table name for DiskLabel
func (DiskLabel) TableName() string {
	return "disk_labels"
}
