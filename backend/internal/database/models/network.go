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
	Status      string    `gorm:"default:pending" json:"status"` // pending, pending_changes, active, error, rollback
	LastError   string    `gorm:"type:text" json:"last_error,omitempty"`

	// Pending changes tracking (Proxmox-style)
	HasPendingChanges bool      `gorm:"default:false" json:"has_pending_changes"` // True if changes not yet applied
	PendingPorts      string    `gorm:"type:text" json:"pending_ports,omitempty"` // Pending ports to apply
	PendingIPAddress  string    `json:"pending_ip_address,omitempty"` // Pending IP to apply
	PendingGateway    string    `json:"pending_gateway,omitempty"` // Pending gateway to apply

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

// NetworkSnapshot represents a snapshot of the network configuration before applying changes
// This enables rollback to the previous working state if changes break connectivity
type NetworkSnapshot struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	BridgeName  string    `gorm:"index" json:"bridge_name"` // Which bridge this snapshot is for

	// Snapshot of current configuration (before applying changes)
	CurrentPorts     string    `gorm:"type:text" json:"current_ports"`
	CurrentIPAddress string    `json:"current_ip_address,omitempty"`
	CurrentGateway   string    `json:"current_gateway,omitempty"`

	// Interface states before changes
	InterfaceStates  string    `gorm:"type:text" json:"interface_states"` // JSON: map[interface]state
	RouteTable       string    `gorm:"type:text" json:"route_table"` // Output of "ip route show"

	// Metadata
	CreatedAt        time.Time `json:"created_at"`
	AppliedAt        time.Time `json:"applied_at,omitempty"` // When changes were applied
	RolledBackAt     time.Time `json:"rolled_back_at,omitempty"` // When rollback occurred
	Status           string    `gorm:"default:active" json:"status"` // active, applied, rolled_back
}
