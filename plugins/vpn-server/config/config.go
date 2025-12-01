package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the VPN server
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	WireGuard WireGuardConfig `mapstructure:"wireguard"`
	OpenVPN   OpenVPNConfig   `mapstructure:"openvpn"`
	PPTP      PPTPConfig      `mapstructure:"pptp"`
	L2TP      L2TPConfig      `mapstructure:"l2tp"`
	General   GeneralConfig   `mapstructure:"general"`
	Security  SecurityConfig  `mapstructure:"security"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Charset  string `mapstructure:"charset"`
}

// WireGuardConfig holds WireGuard server configuration
type WireGuardConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Interface  string `mapstructure:"interface"`
	ListenPort int    `mapstructure:"listenPort"`
	Subnet     string `mapstructure:"subnet"`
	DNS        string `mapstructure:"dns"`
	Endpoint   string `mapstructure:"endpoint"`
	MTU        int    `mapstructure:"mtu"`
}

// OpenVPNConfig holds OpenVPN server configuration
type OpenVPNConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Protocol    string `mapstructure:"protocol"` // udp, tcp
	Port        int    `mapstructure:"port"`
	Subnet      string `mapstructure:"subnet"`
	Cipher      string `mapstructure:"cipher"`
	Auth        string `mapstructure:"auth"`
	Compression string `mapstructure:"compression"`
	TLSVersion  string `mapstructure:"tlsVersion"`
	MaxClients  int    `mapstructure:"maxClients"`
	CAPath      string `mapstructure:"caPath"`
	CertPath    string `mapstructure:"certPath"`
	KeyPath     string `mapstructure:"keyPath"`
	DHPath      string `mapstructure:"dhPath"`
}

// PPTPConfig holds PPTP server configuration
type PPTPConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Port           int    `mapstructure:"port"`
	Subnet         string `mapstructure:"subnet"`
	Encryption     string `mapstructure:"encryption"`
	Authentication string `mapstructure:"authentication"`
	LocalIP        string `mapstructure:"localIP"`
	RemoteIP       string `mapstructure:"remoteIP"`
}

// L2TPConfig holds L2TP/IPsec server configuration
type L2TPConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Port           int    `mapstructure:"port"`
	IPsecPort      int    `mapstructure:"ipsecPort"`
	Subnet         string `mapstructure:"subnet"`
	PSK            string `mapstructure:"psk"`
	Encryption     string `mapstructure:"encryption"`
	Authentication string `mapstructure:"authentication"`
	NATTraversal   bool   `mapstructure:"natTraversal"`
	LocalIP        string `mapstructure:"localIP"`
	IPRange        string `mapstructure:"ipRange"`
}

// GeneralConfig holds general VPN settings
type GeneralConfig struct {
	DefaultInterface         string `mapstructure:"defaultInterface"`
	AccountSource            string `mapstructure:"accountSource"` // local, ldap, ad, radius
	EnableLogging            bool   `mapstructure:"enableLogging"`
	LogLevel                 string `mapstructure:"logLevel"`
	MaxConcurrentConnections int    `mapstructure:"maxConcurrentConnections"`
	ConnectionTimeout        int    `mapstructure:"connectionTimeout"`
	EnableIPv6               bool   `mapstructure:"enableIPv6"`
	DNSServers               string `mapstructure:"dnsServers"`
	DefaultGateway           string `mapstructure:"defaultGateway"`
	EnableNAT                bool   `mapstructure:"enableNAT"`
	ForwardingRules          string `mapstructure:"forwardingRules"`
}

// SecurityConfig holds security-related settings
type SecurityConfig struct {
	JWTSecret           string `mapstructure:"jwtSecret"`
	JWTExpiration       int    `mapstructure:"jwtExpiration"` // hours
	RateLimitPerMinute  int    `mapstructure:"rateLimitPerMinute"`
	EnableAuditLogging  bool   `mapstructure:"enableAuditLogging"`
	AuditLogPath        string `mapstructure:"auditLogPath"`
	AllowedOrigins      []string `mapstructure:"allowedOrigins"`
}

// Load reads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("vpn-server")
		v.SetConfigType("yaml")
		v.AddConfigPath("/etc/stumpfworks-nas/")
		v.AddConfigPath("$HOME/.stumpfworks-nas/")
		v.AddConfigPath(".")
	}

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("VPN")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found; using defaults
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override from environment variables
	if dbPass := os.Getenv("VPN_DATABASE_PASSWORD"); dbPass != "" {
		cfg.Database.Password = dbPass
	}
	if jwtSecret := os.Getenv("VPN_SECURITY_JWTSECRET"); jwtSecret != "" {
		cfg.Security.JWTSecret = jwtSecret
	}
	if l2tpPSK := os.Getenv("VPN_L2TP_PSK"); l2tpPSK != "" {
		cfg.L2TP.PSK = l2tpPSK
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "release")

	// Database defaults
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.database", "stumpfworks_nas")
	v.SetDefault("database.charset", "utf8mb4")

	// WireGuard defaults
	v.SetDefault("wireguard.enabled", true)
	v.SetDefault("wireguard.interface", "wg0")
	v.SetDefault("wireguard.listenPort", 51820)
	v.SetDefault("wireguard.subnet", "10.8.0.0/24")
	v.SetDefault("wireguard.dns", "8.8.8.8, 8.8.4.4")
	v.SetDefault("wireguard.mtu", 1420)

	// OpenVPN defaults
	v.SetDefault("openvpn.enabled", true)
	v.SetDefault("openvpn.protocol", "udp")
	v.SetDefault("openvpn.port", 1194)
	v.SetDefault("openvpn.subnet", "10.9.0.0/24")
	v.SetDefault("openvpn.cipher", "AES-256-GCM")
	v.SetDefault("openvpn.auth", "SHA512")
	v.SetDefault("openvpn.compression", "lz4-v2")
	v.SetDefault("openvpn.tlsVersion", "1.3")
	v.SetDefault("openvpn.maxClients", 100)

	// PPTP defaults
	v.SetDefault("pptp.enabled", false)
	v.SetDefault("pptp.port", 1723)
	v.SetDefault("pptp.subnet", "10.10.0.0/24")
	v.SetDefault("pptp.encryption", "MPPE-128")
	v.SetDefault("pptp.authentication", "MS-CHAPv2")

	// L2TP defaults
	v.SetDefault("l2tp.enabled", true)
	v.SetDefault("l2tp.port", 1701)
	v.SetDefault("l2tp.ipsecPort", 500)
	v.SetDefault("l2tp.subnet", "10.11.0.0/24")
	v.SetDefault("l2tp.encryption", "AES-256")
	v.SetDefault("l2tp.authentication", "SHA2-256")
	v.SetDefault("l2tp.natTraversal", true)

	// General defaults
	v.SetDefault("general.accountSource", "local")
	v.SetDefault("general.enableLogging", true)
	v.SetDefault("general.logLevel", "info")
	v.SetDefault("general.maxConcurrentConnections", 100)
	v.SetDefault("general.connectionTimeout", 300)
	v.SetDefault("general.enableIPv6", false)
	v.SetDefault("general.dnsServers", "8.8.8.8, 8.8.4.4")
	v.SetDefault("general.enableNAT", true)

	// Security defaults
	v.SetDefault("security.jwtExpiration", 24)
	v.SetDefault("security.rateLimitPerMinute", 60)
	v.SetDefault("security.enableAuditLogging", true)
	v.SetDefault("security.auditLogPath", "/var/log/vpn-server/audit.log")
	v.SetDefault("security.allowedOrigins", []string{"*"})
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	if c.L2TP.Enabled && c.L2TP.PSK == "" {
		return fmt.Errorf("L2TP PSK is required when L2TP is enabled")
	}
	return nil
}
