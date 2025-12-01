// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package models

import (
	"time"
)

// NetworkBridge represents a network bridge configuration in the database
// This ensures bridges persist across reboots and are automatically restored on startup
type NetworkBridge struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"` // e.g., br0, vmbr0, vmbr1
	Description string    `json:"description"`
	Ports       string    `gorm:"type:text" json:"ports"` // Comma-separated list of interface names (e.g., "ens18,ens19")
	IPAddress   string    `json:"ip_address,omitempty"` // Optional static IP (CIDR notation: 192.168.1.10/24)
	Gateway     string    `json:"gateway,omitempty"`    // Optional gateway
	Autostart   bool      `gorm:"default:true" json:"autostart"` // Auto-create on system boot
	Status      string    `gorm:"default:pending" json:"status"` // pending, active, error
	LastError   string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for the NetworkBridge model
func (NetworkBridge) TableName() string {
	return "network_bridges"
}

// NetworkInterface represents a network interface configuration
type NetworkInterface struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"` // e.g., ens18, eth0
	Method    string    `gorm:"default:dhcp" json:"method"` // dhcp, static, manual
	IPAddress string    `json:"ip_address,omitempty"` // CIDR notation
	Gateway   string    `json:"gateway,omitempty"`
	DNS       string    `gorm:"type:text" json:"dns,omitempty"` // Comma-separated DNS servers
	MTU       int       `json:"mtu,omitempty"`
	Autostart bool      `gorm:"default:true" json:"autostart"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for the NetworkInterface model
func (NetworkInterface) TableName() string {
	return "network_interfaces"
}
