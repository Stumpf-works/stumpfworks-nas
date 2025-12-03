package models

import (
	"time"

	"gorm.io/gorm"
)

// TimeMachineDevice represents a macOS device backing up via Time Machine
type TimeMachineDevice struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Device information
	DeviceName  string    `gorm:"uniqueIndex;not null" json:"device_name"` // Mac hostname
	MACAddress  string    `gorm:"index" json:"mac_address"`                // Network MAC address for identification
	ModelID     string    `json:"model_id"`                                 // Mac model identifier (e.g., "MacBookPro18,1")

	// Storage configuration
	SharePath   string    `gorm:"not null" json:"share_path"`              // Path where backups are stored
	QuotaGB     int       `gorm:"default:0" json:"quota_gb"`               // Quota in GB (0 = unlimited)
	UsedGB      float64   `json:"used_gb"`                                  // Currently used space in GB

	// Status
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	LastBackup  *time.Time `json:"last_backup"`                             // Last successful backup timestamp
	LastSeen    *time.Time `json:"last_seen"`                               // Last time device connected

	// Authentication
	Username    string    `json:"username"`                                 // SMB username for this device
	Password    string    `json:"-"`                                        // SMB password (never sent to frontend)
}

// TimeMachineConfig stores global Time Machine server configuration
type TimeMachineConfig struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Service settings
	Enabled         bool      `gorm:"default:false" json:"enabled"`
	ShareName       string    `gorm:"default:'TimeMachine'" json:"share_name"`
	BasePath        string    `gorm:"default:'/mnt/storage/timemachine'" json:"base_path"`

	// Default settings for new devices
	DefaultQuotaGB  int       `gorm:"default:500" json:"default_quota_gb"`
	AutoDiscovery   bool      `gorm:"default:true" json:"auto_discovery"`   // Advertise via Avahi/Bonjour

	// Advanced settings
	UseAFP          bool      `gorm:"default:false" json:"use_afp"`         // Use AFP instead of SMB (legacy)
	UseSMB          bool      `gorm:"default:true" json:"use_smb"`          // Use SMB (recommended)
	SMBVersion      string    `gorm:"default:'3'" json:"smb_version"`       // SMB protocol version
}
