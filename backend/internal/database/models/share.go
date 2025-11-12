package models

import (
	"gorm.io/gorm"
)

// Share represents a network share in the database
type Share struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;size:255;not null"`
	Path        string `gorm:"size:500;not null"`
	Type        string `gorm:"size:10;not null"` // smb, nfs, ftp
	Description string `gorm:"size:500"`
	Enabled     bool   `gorm:"default:true"`
	ReadOnly    bool   `gorm:"default:false"`
	Browseable  bool   `gorm:"default:true"`
	GuestOK     bool   `gorm:"default:false"`
	ValidUsers  string `gorm:"size:1000"` // Comma-separated list
}

// TableName specifies the table name for Share
func (Share) TableName() string {
	return "shares"
}
