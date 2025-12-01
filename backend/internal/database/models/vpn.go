// Revision: 2025-12-01 | Author: Claude | Version: 1.0.0
package models

import (
	"time"
)

// VPNProtocolConfig represents VPN protocol configuration stored in the database
// This ensures protocol settings persist across reboots
type VPNProtocolConfig struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Protocol      string    `gorm:"uniqueIndex;not null" json:"protocol"` // wireguard, openvpn, pptp, l2tp
	Enabled       bool      `gorm:"default:false" json:"enabled"`
	ListenPort    int       `json:"listen_port"`
	NetworkRange  string    `json:"network_range"`  // e.g., "10.8.0.0/24"
	DNS           string    `json:"dns"`            // Comma-separated DNS servers
	PublicKey     string    `gorm:"type:text" json:"public_key,omitempty"`
	PrivateKey    string    `gorm:"type:text" json:"private_key,omitempty"` // Encrypted in production
	ServerAddress string    `json:"server_address"` // Public IP/hostname for client configs
	ConfigPath    string    `json:"config_path"`    // Path to config file
	ExtraConfig   string    `gorm:"type:text" json:"extra_config,omitempty"` // JSON blob for protocol-specific settings
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name for the VPNProtocolConfig model
func (VPNProtocolConfig) TableName() string {
	return "vpn_protocol_configs"
}

// VPNPeer represents a WireGuard peer stored in the database
type VPNPeer struct {
	ID              string     `gorm:"primaryKey" json:"id"`
	ProtocolID      uint       `gorm:"not null;index" json:"protocol_id"` // Foreign key to VPNProtocolConfig
	Name            string     `gorm:"not null" json:"name"`
	Description     string     `gorm:"type:text" json:"description,omitempty"`
	PublicKey       string     `gorm:"uniqueIndex;not null;type:text" json:"public_key"`
	PrivateKey      string     `gorm:"type:text" json:"private_key,omitempty"` // Only stored for server-generated keys
	PresharedKey    string     `gorm:"type:text" json:"preshared_key,omitempty"`
	AllowedIPs      string     `gorm:"not null" json:"allowed_ips"` // e.g., "10.8.0.2/32"
	Endpoint        string     `json:"endpoint,omitempty"`          // For site-to-site configs
	PersistentKeepalive int    `gorm:"default:0" json:"persistent_keepalive"`
	Enabled         bool       `gorm:"default:true" json:"enabled"`
	BytesReceived   uint64     `gorm:"default:0" json:"bytes_received"`
	BytesSent       uint64     `gorm:"default:0" json:"bytes_sent"`
	LatestHandshake *time.Time `json:"latest_handshake,omitempty"`
	LastSeen        *time.Time `json:"last_seen,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationship
	Protocol VPNProtocolConfig `gorm:"foreignKey:ProtocolID" json:"protocol,omitempty"`
}

// TableName specifies the table name for the VPNPeer model
func (VPNPeer) TableName() string {
	return "vpn_peers"
}

// VPNCertificate represents an OpenVPN certificate stored in the database
type VPNCertificate struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	ProtocolID   uint       `gorm:"not null;index" json:"protocol_id"` // Foreign key to VPNProtocolConfig
	CommonName   string     `gorm:"not null;uniqueIndex:idx_protocol_cn" json:"common_name"`
	Description  string     `gorm:"type:text" json:"description,omitempty"`
	SerialNumber string     `gorm:"uniqueIndex" json:"serial_number"`
	CertData     string     `gorm:"type:text" json:"cert_data,omitempty"` // PEM encoded certificate
	KeyData      string     `gorm:"type:text" json:"key_data,omitempty"`  // PEM encoded private key (encrypted)
	ValidFrom    time.Time  `json:"valid_from"`
	ValidTo      time.Time  `json:"valid_to"`
	Status       string     `gorm:"default:valid" json:"status"` // valid, revoked, expired
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	RevokeReason string     `json:"revoke_reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relationship
	Protocol VPNProtocolConfig `gorm:"foreignKey:ProtocolID" json:"protocol,omitempty"`
}

// TableName specifies the table name for the VPNCertificate model
func (VPNCertificate) TableName() string {
	return "vpn_certificates"
}

// VPNUser represents a VPN user account (cross-protocol)
// This allows unified user management across all VPN protocols
type VPNUser struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Username          string    `gorm:"uniqueIndex;not null" json:"username"`
	Email             string    `json:"email,omitempty"`
	PasswordHash      string    `gorm:"not null" json:"-"` // For PPTP/L2TP/OpenVPN username/password auth
	Enabled           bool      `gorm:"default:true" json:"enabled"`
	MaxConnections    int       `gorm:"default:1" json:"max_connections"` // Max simultaneous connections
	AllowedProtocols  string    `gorm:"type:text" json:"allowed_protocols"` // Comma-separated: "wireguard,openvpn"
	IPRestrictions    string    `gorm:"type:text" json:"ip_restrictions,omitempty"` // JSON array of allowed source IPs
	BandwidthLimit    uint64    `gorm:"default:0" json:"bandwidth_limit"` // Bytes per second, 0 = unlimited
	DataQuota         uint64    `gorm:"default:0" json:"data_quota"` // Total bytes allowed, 0 = unlimited
	DataUsed          uint64    `gorm:"default:0" json:"data_used"` // Bytes used (reset monthly)
	QuotaResetDay     int       `gorm:"default:1" json:"quota_reset_day"` // Day of month to reset quota
	LastConnection    *time.Time `json:"last_connection,omitempty"`
	LastDisconnection *time.Time `json:"last_disconnection,omitempty"`
	ConnectionCount   uint64    `gorm:"default:0" json:"connection_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName specifies the table name for the VPNUser model
func (VPNUser) TableName() string {
	return "vpn_users"
}

// VPNConnection represents an active or historical VPN connection
// This enables connection monitoring and audit logging
type VPNConnection struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"index" json:"user_id,omitempty"`
	PeerID         string     `gorm:"index" json:"peer_id,omitempty"` // For WireGuard
	CertificateID  string     `gorm:"index" json:"certificate_id,omitempty"` // For OpenVPN
	Protocol       string     `gorm:"not null;index" json:"protocol"`
	ClientIP       string     `json:"client_ip"`        // VPN tunnel IP assigned to client
	SourceIP       string     `json:"source_ip"`        // Real client IP
	VirtualIP      string     `json:"virtual_ip"`       // VPN IP assigned
	BytesReceived  uint64     `gorm:"default:0" json:"bytes_received"`
	BytesSent      uint64     `gorm:"default:0" json:"bytes_sent"`
	ConnectedAt    time.Time  `gorm:"index" json:"connected_at"`
	DisconnectedAt *time.Time `gorm:"index" json:"disconnected_at,omitempty"`
	Duration       int        `json:"duration"` // Seconds
	DisconnectReason string   `json:"disconnect_reason,omitempty"`
	Active         bool       `gorm:"default:true;index" json:"active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	User        *VPNUser        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Peer        *VPNPeer        `gorm:"foreignKey:PeerID" json:"peer,omitempty"`
	Certificate *VPNCertificate `gorm:"foreignKey:CertificateID" json:"certificate,omitempty"`
}

// TableName specifies the table name for the VPNConnection model
func (VPNConnection) TableName() string {
	return "vpn_connections"
}

// VPNRoute represents static routes pushed to VPN clients
type VPNRoute struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProtocolID  uint      `gorm:"not null;index" json:"protocol_id"`
	Network     string    `gorm:"not null" json:"network"` // e.g., "192.168.1.0/24"
	Gateway     string    `json:"gateway,omitempty"`
	Metric      int       `gorm:"default:0" json:"metric"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationship
	Protocol VPNProtocolConfig `gorm:"foreignKey:ProtocolID" json:"protocol,omitempty"`
}

// TableName specifies the table name for the VPNRoute model
func (VPNRoute) TableName() string {
	return "vpn_routes"
}

// VPNFirewallRule represents firewall rules for VPN traffic
type VPNFirewallRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProtocolID  uint      `gorm:"not null;index" json:"protocol_id"`
	Name        string    `gorm:"not null" json:"name"`
	Action      string    `gorm:"not null" json:"action"` // allow, deny
	Direction   string    `gorm:"not null" json:"direction"` // inbound, outbound, forward
	SourceIP    string    `json:"source_ip,omitempty"`    // CIDR or "any"
	DestIP      string    `json:"dest_ip,omitempty"`      // CIDR or "any"
	SourcePort  string    `json:"source_port,omitempty"`  // Port or range
	DestPort    string    `json:"dest_port,omitempty"`    // Port or range
	Protocol    string    `json:"protocol,omitempty"`     // tcp, udp, icmp, any
	Priority    int       `gorm:"default:100" json:"priority"` // Lower = higher priority
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationship
	VPNProtocol VPNProtocolConfig `gorm:"foreignKey:ProtocolID" json:"vpn_protocol,omitempty"`
}

// TableName specifies the table name for the VPNFirewallRule model
func (VPNFirewallRule) TableName() string {
	return "vpn_firewall_rules"
}
