// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package models

import (
	"time"
)

// LXCContainer represents an LXC container configuration in the database
// This ensures containers persist across reboots and maintain their configuration
type LXCContainer struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"uniqueIndex;not null" json:"name"`
	Template     string    `gorm:"not null" json:"template"` // ubuntu, debian, alpine, centos
	Release      string    `gorm:"not null" json:"release"`  // 22.04, bullseye, 3.18, etc.
	Architecture string    `gorm:"not null" json:"architecture"` // amd64, arm64
	MemoryLimit  int64     `json:"memory_limit,omitempty"` // MB
	CPULimit     int       `json:"cpu_limit,omitempty"` // Number of CPUs
	Autostart    bool      `gorm:"default:false" json:"autostart"` // Auto-start on boot
	NetworkMode  string    `gorm:"default:internal" json:"network_mode"` // internal (lxcbr0), bridged
	Bridge       string    `json:"bridge,omitempty"` // Bridge name when network_mode is "bridged"
	IPv4         string    `json:"ipv4,omitempty"` // Current IPv4 address (updated dynamically)
	IPv6         string    `json:"ipv6,omitempty"` // Current IPv6 address (updated dynamically)
	Status       string    `gorm:"default:stopped" json:"status"` // stopped, running, frozen, error
	LastError    string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the table name for the LXCContainer model
func (LXCContainer) TableName() string {
	return "lxc_containers"
}
