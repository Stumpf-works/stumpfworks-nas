package database

import (
	"time"

	"gorm.io/gorm"
)

// VPNUser represents a VPN user with protocol access permissions
type VPNUser struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Username       string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
	Email          string    `gorm:"type:varchar(255);not null" json:"email"`
	PasswordHash   string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	LastConnection *time.Time `json:"lastConnection,omitempty"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`

	// Relationships
	Protocols  []VPNUserProtocol  `gorm:"foreignKey:UserID" json:"protocols,omitempty"`
	Peers      []WireGuardPeer    `gorm:"foreignKey:UserID" json:"peers,omitempty"`
	Certificates []OpenVPNCertificate `gorm:"foreignKey:UserID" json:"certificates,omitempty"`
	Connections []VPNConnection    `gorm:"foreignKey:UserID" json:"connections,omitempty"`
}

// VPNUserProtocol represents protocol access permissions for a user
type VPNUserProtocol struct {
	UserID   string `gorm:"type:varchar(36);primaryKey" json:"userId"`
	Protocol string `gorm:"type:varchar(20);primaryKey" json:"protocol"` // wireguard, openvpn, pptp, l2tp
	Enabled  bool   `gorm:"default:false" json:"enabled"`

	User VPNUser `gorm:"foreignKey:UserID" json:"-"`
}

// WireGuardPeer represents a WireGuard peer configuration
type WireGuardPeer struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID     *string   `gorm:"type:varchar(36);index" json:"userId,omitempty"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	PublicKey  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"publicKey"`
	PrivateKey string    `gorm:"type:varchar(255);not null" json:"privateKey,omitempty"`
	AllowedIPs string    `gorm:"type:varchar(255);not null" json:"allowedIPs"`
	Endpoint   string    `gorm:"type:varchar(255)" json:"endpoint,omitempty"`
	Enabled    bool      `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time `json:"createdAt"`

	User *VPNUser `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OpenVPNCertificate represents an OpenVPN client certificate
type OpenVPNCertificate struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID       *string   `gorm:"type:varchar(36);index" json:"userId,omitempty"`
	CommonName   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"commonName"`
	SerialNumber string    `gorm:"type:varchar(255);not null" json:"serialNumber"`
	ValidFrom    time.Time `json:"validFrom"`
	ValidTo      time.Time `json:"validTo"`
	Status       string    `gorm:"type:varchar(20);default:'valid'" json:"status"` // valid, revoked, expired
	Certificate  string    `gorm:"type:text;not null" json:"certificate"`
	PrivateKey   string    `gorm:"type:text;not null" json:"privateKey,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`

	User *VPNUser `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// VPNConnection represents an active or historical VPN connection
type VPNConnection struct {
	ID             string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID         string     `gorm:"type:varchar(36);not null;index:idx_user_protocol" json:"userId"`
	Protocol       string     `gorm:"type:varchar(20);not null;index:idx_user_protocol" json:"protocol"`
	IPAddress      string     `gorm:"type:varchar(45);not null" json:"ipAddress"`
	ConnectedAt    time.Time  `gorm:"index:idx_connected" json:"connectedAt"`
	DisconnectedAt *time.Time `json:"disconnectedAt,omitempty"`
	BytesReceived  int64      `gorm:"default:0" json:"bytesReceived"`
	BytesSent      int64      `gorm:"default:0" json:"bytesSent"`

	User VPNUser `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// VPNConfig represents server configuration for each protocol
type VPNConfig struct {
	Protocol  string    `gorm:"type:varchar(20);primaryKey" json:"protocol"`
	Config    string    `gorm:"type:json;not null" json:"config"`
	Enabled   bool      `gorm:"default:false" json:"enabled"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// WireGuardConfig represents WireGuard server configuration
type WireGuardConfig struct {
	Interface    string `json:"interface"`
	PrivateKey   string `json:"privateKey"`
	PublicKey    string `json:"publicKey"`
	ListenPort   int    `json:"listenPort"`
	Endpoint     string `json:"endpoint"`
	Subnet       string `json:"subnet"`
	DNS          string `json:"dns"`
	PostUp       string `json:"postUp,omitempty"`
	PostDown     string `json:"postDown,omitempty"`
	MTU          int    `json:"mtu,omitempty"`
}

// OpenVPNConfig represents OpenVPN server configuration
type OpenVPNConfig struct {
	Protocol       string `json:"protocol"` // UDP or TCP
	Port           int    `json:"port"`
	Subnet         string `json:"subnet"`
	Cipher         string `json:"cipher"`
	Auth           string `json:"auth"`
	Compression    string `json:"compression"`
	TLSVersion     string `json:"tlsVersion"`
	MaxClients     int    `json:"maxClients"`
	KeepAlive      string `json:"keepAlive"`
	VerbosityLevel int    `json:"verbosityLevel"`
}

// PPTPConfig represents PPTP server configuration
type PPTPConfig struct {
	Port           int    `json:"port"`
	Subnet         string `json:"subnet"`
	Encryption     string `json:"encryption"` // MPPE-128, MPPE-40
	Authentication string `json:"authentication"` // MS-CHAPv2, CHAP, PAP
	RequireMPPE    bool   `json:"requireMPPE"`
	LocalIP        string `json:"localIP"`
	RemoteIP       string `json:"remoteIP"`
}

// L2TPConfig represents L2TP/IPsec server configuration
type L2TPConfig struct {
	Port           int    `json:"port"`
	IPsecPort      int    `json:"ipsecPort"`
	Subnet         string `json:"subnet"`
	PSK            string `json:"psk"`
	Encryption     string `json:"encryption"` // AES-256, AES-192, AES-128, 3DES
	Authentication string `json:"authentication"` // SHA2-256, SHA2-512, SHA1
	NATTraversal   bool   `json:"natTraversal"`
	LocalIP        string `json:"localIP"`
	IPRange        string `json:"ipRange"`
}

// TableName overrides for GORM
func (VPNUser) TableName() string {
	return "vpn_users"
}

func (VPNUserProtocol) TableName() string {
	return "vpn_user_protocols"
}

func (WireGuardPeer) TableName() string {
	return "wireguard_peers"
}

func (OpenVPNCertificate) TableName() string {
	return "openvpn_certificates"
}

func (VPNConnection) TableName() string {
	return "vpn_connections"
}

func (VPNConfig) TableName() string {
	return "vpn_config"
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&VPNUser{},
		&VPNUserProtocol{},
		&WireGuardPeer{},
		&OpenVPNCertificate{},
		&VPNConnection{},
		&VPNConfig{},
	)
}
