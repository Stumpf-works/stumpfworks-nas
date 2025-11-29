package models

import (
	"time"

	"gorm.io/gorm"
)

// AddonInstallation tracks installed addons
type AddonInstallation struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AddonID     string         `gorm:"uniqueIndex;not null" json:"addon_id"` // e.g., "vm-manager"
	Version     string         `json:"version"`                               // Installed version
	Installed   bool           `gorm:"default:false" json:"installed"`
	InstallDate time.Time      `json:"install_date"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Error       string         `gorm:"type:text" json:"error,omitempty"` // Installation error if any
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for AddonInstallation
func (AddonInstallation) TableName() string {
	return "addon_installations"
}
