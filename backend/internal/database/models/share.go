// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package models

import (
	"gorm.io/gorm"
)

// Share represents a network share in the database
type Share struct {
	gorm.Model
	Name        string `gorm:"size:255;not null;uniqueIndex:idx_name_deleted"` // Composite unique with deleted_at
	Path        string `gorm:"size:500;not null"`
	VolumeID    string `gorm:"size:100;index"` // Optional - links to a managed volume
	Type        string `gorm:"size:10;not null"` // smb, nfs, ftp
	Description string `gorm:"size:500"`
	Enabled     bool   `gorm:"default:true"`
	ReadOnly    bool   `gorm:"default:false"`
	Browseable  bool   `gorm:"default:true"`
	GuestOK     bool   `gorm:"default:false"`
	ValidUsers  string `gorm:"size:1000"` // Comma-separated list of usernames
	ValidGroups string `gorm:"size:1000"` // Comma-separated list of group names
	DeletedAt   gorm.DeletedAt `gorm:"index;uniqueIndex:idx_name_deleted"` // Part of composite unique index
}

// TableName specifies the table name for Share
func (Share) TableName() string {
	return "shares"
}
