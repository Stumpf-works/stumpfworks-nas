// Package vpn provides VPN server management with multi-protocol support
package vpn

import (
	"fmt"
	"strings"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/executor"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// VPNManager manages multiple VPN protocols with database persistence
type VPNManager struct {
	shell executor.ShellExecutor
	db    *gorm.DB
}

// Protocol represents a VPN protocol configuration
type Protocol string

const (
	WireGuard Protocol = "wireguard"
	OpenVPN   Protocol = "openvpn"
	PPTP      Protocol = "pptp"
	L2TP      Protocol = "l2tp"
)

// ProtocolStatus represents the status of a VPN protocol
type ProtocolStatus struct {
	Protocol    string `json:"protocol"`
	Installed   bool   `json:"installed"`
	Enabled     bool   `json:"enabled"`
	Running     bool   `json:"running"`
	Connections int    `json:"connections"`
	Error       string `json:"error,omitempty"`
}

// ProtocolPackages maps protocols to their required system packages
var ProtocolPackages = map[Protocol][]string{
	WireGuard: {
		"wireguard",
		"wireguard-tools",
		"qrencode", // For QR code generation
	},
	OpenVPN: {
		"openvpn",
		"easy-rsa", // For certificate management
	},
	PPTP: {
		"pptpd",
	},
	L2TP: {
		"xl2tpd",
		"strongswan", // For IPsec
	},
}

// ProtocolServices maps protocols to their systemd services
var ProtocolServices = map[Protocol][]string{
	WireGuard: {"wg-quick@wg0"}, // WireGuard interface service
	OpenVPN:   {"openvpn@server"},
	PPTP:      {"pptpd"},
	L2TP:      {"xl2tpd", "strongswan"},
}

// NewVPNManager creates a new VPN manager with database persistence
func NewVPNManager(shell executor.ShellExecutor) *VPNManager {
	return &VPNManager{
		shell: shell,
		db:    database.GetDB(),
	}
}

// GetProtocolStatus returns the status of a specific protocol
func (vm *VPNManager) GetProtocolStatus(protocol Protocol) (*ProtocolStatus, error) {
	status := &ProtocolStatus{
		Protocol:  string(protocol),
		Installed: false,
		Enabled:   false,
		Running:   false,
	}

	// Check if packages are installed
	packages := ProtocolPackages[protocol]
	if len(packages) > 0 {
		allInstalled := true
		for _, pkg := range packages {
			result, err := vm.shell.Execute("dpkg", "-l", pkg)
			if err != nil || !strings.Contains(result.Stdout, "ii") {
				allInstalled = false
				break
			}
		}
		status.Installed = allInstalled
	}

	// If not installed, return early
	if !status.Installed {
		return status, nil
	}

	// Check if services are enabled and running
	services := ProtocolServices[protocol]
	if len(services) > 0 {
		for _, service := range services {
			// Check if enabled
			result, err := vm.shell.Execute("systemctl", "is-enabled", service)
			if err == nil && strings.TrimSpace(result.Stdout) == "enabled" {
				status.Enabled = true
			}

			// Check if running
			result, err = vm.shell.Execute("systemctl", "is-active", service)
			if err == nil && strings.TrimSpace(result.Stdout) == "active" {
				status.Running = true
			}
		}
	}

	// Get connection count from database
	if vm.db != nil {
		connections, err := vm.GetConnectionsByProtocol(protocol)
		if err == nil {
			status.Connections = len(connections)
		}
	}

	return status, nil
}

// GetAllProtocolStatuses returns status of all protocols
func (vm *VPNManager) GetAllProtocolStatuses() ([]ProtocolStatus, error) {
	protocols := []Protocol{WireGuard, OpenVPN, PPTP, L2TP}
	statuses := make([]ProtocolStatus, 0, len(protocols))

	for _, protocol := range protocols {
		status, err := vm.GetProtocolStatus(protocol)
		if err != nil {
			logger.Warn("Failed to get protocol status",
				zap.String("protocol", string(protocol)),
				zap.Error(err))
			continue
		}
		statuses = append(statuses, *status)
	}

	return statuses, nil
}

// InstallProtocol installs packages for a specific VPN protocol
func (vm *VPNManager) InstallProtocol(protocol Protocol) error {
	logger.Info("Installing VPN protocol", zap.String("protocol", string(protocol)))

	packages := ProtocolPackages[protocol]
	if len(packages) == 0 {
		return fmt.Errorf("no packages defined for protocol: %s", protocol)
	}

	// Update package lists
	logger.Info("Updating package lists...")
	if _, err := vm.shell.Execute("apt-get", "update"); err != nil {
		return fmt.Errorf("failed to update package lists: %w", err)
	}

	// Install packages
	args := append([]string{"install", "-y"}, packages...)
	logger.Info("Installing packages", zap.Strings("packages", packages))

	// Use longer timeout for package installation
	result, err := vm.shell.ExecuteWithTimeout(10*time.Minute, "apt-get", args...)
	if err != nil {
		return fmt.Errorf("failed to install packages: %s: %w", result.Stderr, err)
	}

	logger.Info("Protocol packages installed successfully",
		zap.String("protocol", string(protocol)))

	return nil
}

// InitializeProtocol initializes a protocol after installation (creates config files, etc.)
func (vm *VPNManager) InitializeProtocol(protocol Protocol) error {
	logger.Info("Initializing VPN protocol", zap.String("protocol", string(protocol)))

	// Create or get protocol config from database
	config, err := vm.GetOrCreateProtocolConfig(protocol)
	if err != nil {
		logger.Warn("Failed to create protocol config in database",
			zap.String("protocol", string(protocol)),
			zap.Error(err))
	}

	// Initialize protocol-specific configuration
	var initErr error
	switch protocol {
	case WireGuard:
		initErr = vm.initializeWireGuard()
		if initErr == nil && config != nil {
			// Extract and save public key from initialization
			result, _ := vm.shell.Execute("cat", "/etc/wireguard/wg0.conf")
			if strings.Contains(result.Stdout, "PrivateKey") {
				lines := strings.Split(result.Stdout, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "PrivateKey") {
						config.PrivateKey = strings.TrimSpace(strings.TrimPrefix(line, "PrivateKey = "))
					}
				}
				vm.SaveProtocolConfig(config)
			}
		}
	case OpenVPN:
		initErr = vm.initializeOpenVPN()
	case PPTP:
		initErr = vm.initializePPTP()
	case L2TP:
		initErr = vm.initializeL2TP()
	default:
		return fmt.Errorf("unknown protocol: %s", protocol)
	}

	return initErr
}

// EnableProtocol enables and starts services for a protocol
func (vm *VPNManager) EnableProtocol(protocol Protocol) error {
	logger.Info("Enabling VPN protocol", zap.String("protocol", string(protocol)))

	services := ProtocolServices[protocol]
	for _, service := range services {
		// Enable service
		if _, err := vm.shell.Execute("systemctl", "enable", service); err != nil {
			logger.Warn("Failed to enable service",
				zap.String("service", service),
				zap.Error(err))
		}

		// Start service
		if _, err := vm.shell.Execute("systemctl", "start", service); err != nil {
			return fmt.Errorf("failed to start service %s: %w", service, err)
		}
	}

	// Mark as enabled in database
	if err := vm.MarkProtocolEnabled(protocol, true); err != nil {
		logger.Warn("Failed to mark protocol as enabled in database",
			zap.String("protocol", string(protocol)),
			zap.Error(err))
	}

	logger.Info("Protocol enabled successfully", zap.String("protocol", string(protocol)))
	return nil
}

// DisableProtocol stops and disables services for a protocol
func (vm *VPNManager) DisableProtocol(protocol Protocol) error {
	logger.Info("Disabling VPN protocol", zap.String("protocol", string(protocol)))

	services := ProtocolServices[protocol]
	for _, service := range services {
		// Stop service
		if _, err := vm.shell.Execute("systemctl", "stop", service); err != nil {
			logger.Warn("Failed to stop service",
				zap.String("service", service),
				zap.Error(err))
		}

		// Disable service
		if _, err := vm.shell.Execute("systemctl", "disable", service); err != nil {
			logger.Warn("Failed to disable service",
				zap.String("service", service),
				zap.Error(err))
		}
	}

	// Mark as disabled in database
	if err := vm.MarkProtocolEnabled(protocol, false); err != nil {
		logger.Warn("Failed to mark protocol as disabled in database",
			zap.String("protocol", string(protocol)),
			zap.Error(err))
	}

	logger.Info("Protocol disabled successfully", zap.String("protocol", string(protocol)))
	return nil
}

// Protocol-specific initialization methods

func (vm *VPNManager) initializeWireGuard() error {
	// Generate server keys
	result, err := vm.shell.Execute("wg", "genkey")
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard private key: %w", err)
	}
	privateKey := strings.TrimSpace(result.Stdout)

	// Generate public key from private key using pipe
	pubKeyCmd := fmt.Sprintf("echo '%s' | wg pubkey", privateKey)
	result, err = vm.shell.Execute("sh", "-c", pubKeyCmd)
	if err != nil {
		return fmt.Errorf("failed to generate WireGuard public key: %w", err)
	}
	publicKey := strings.TrimSpace(result.Stdout)

	// Create WireGuard config directory if it doesn't exist
	vm.shell.Execute("mkdir", "-p", "/etc/wireguard")

	// Create basic server configuration
	config := fmt.Sprintf(`[Interface]
Address = 10.8.0.1/24
ListenPort = 51820
PrivateKey = %s
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
`, privateKey)

	// Write config file
	vm.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s' > /etc/wireguard/wg0.conf", config))
	vm.shell.Execute("chmod", "600", "/etc/wireguard/wg0.conf")

	// Enable IP forwarding
	vm.shell.Execute("sysctl", "-w", "net.ipv4.ip_forward=1")
	vm.shell.Execute("sh", "-c", "echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf")

	logger.Info("WireGuard initialized successfully",
		zap.String("public_key", publicKey))

	return nil
}

func (vm *VPNManager) initializeOpenVPN() error {
	// Create OpenVPN directory structure
	vm.shell.Execute("mkdir", "-p", "/etc/openvpn/server")
	vm.shell.Execute("mkdir", "-p", "/etc/openvpn/client")

	// Initialize PKI with easy-rsa
	vm.shell.Execute("make-cadir", "/etc/openvpn/easy-rsa")

	logger.Info("OpenVPN initialized successfully")
	return nil
}

func (vm *VPNManager) initializePPTP() error {
	// Create PPTP config
	config := `option /etc/ppp/pptpd-options
logwtmp
localip 10.10.0.1
remoteip 10.10.0.100-200
`
	vm.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s' > /etc/pptpd.conf", config))

	logger.Info("PPTP initialized successfully")
	return nil
}

func (vm *VPNManager) initializeL2TP() error {
	// Create xl2tpd config
	xl2tpdConfig := `[global]
port = 1701

[lns default]
ip range = 10.11.0.100-10.11.0.200
local ip = 10.11.0.1
require chap = yes
refuse pap = yes
require authentication = yes
name = L2TPServer
pppoptfile = /etc/ppp/options.xl2tpd
length bit = yes
`
	vm.shell.Execute("sh", "-c", fmt.Sprintf("echo '%s' > /etc/xl2tpd/xl2tpd.conf", xl2tpdConfig))

	logger.Info("L2TP/IPsec initialized successfully")
	return nil
}

// Database persistence methods

// SaveProtocolConfig saves or updates a protocol configuration in the database
func (vm *VPNManager) SaveProtocolConfig(config *models.VPNProtocolConfig) error {
	if vm.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Check if config exists
	var existing models.VPNProtocolConfig
	result := vm.db.Where("protocol = ?", config.Protocol).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new
		return vm.db.Create(config).Error
	} else if result.Error != nil {
		return result.Error
	}

	// Update existing
	config.ID = existing.ID
	return vm.db.Save(config).Error
}

// LoadProtocolConfig loads a protocol configuration from the database
func (vm *VPNManager) LoadProtocolConfig(protocol Protocol) (*models.VPNProtocolConfig, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var config models.VPNProtocolConfig
	if err := vm.db.Where("protocol = ?", string(protocol)).First(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

// GetOrCreateProtocolConfig gets an existing protocol config or creates a default one
func (vm *VPNManager) GetOrCreateProtocolConfig(protocol Protocol) (*models.VPNProtocolConfig, error) {
	config, err := vm.LoadProtocolConfig(protocol)
	if err == nil {
		return config, nil
	}

	// Create default config
	defaultConfig := vm.createDefaultProtocolConfig(protocol)
	if err := vm.SaveProtocolConfig(defaultConfig); err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

// createDefaultProtocolConfig creates a default configuration for a protocol
func (vm *VPNManager) createDefaultProtocolConfig(protocol Protocol) *models.VPNProtocolConfig {
	config := &models.VPNProtocolConfig{
		Protocol: string(protocol),
		Enabled:  false,
	}

	switch protocol {
	case WireGuard:
		config.ListenPort = 51820
		config.NetworkRange = "10.8.0.0/24"
		config.DNS = "1.1.1.1,1.0.0.1"
		config.ConfigPath = "/etc/wireguard/wg0.conf"
	case OpenVPN:
		config.ListenPort = 1194
		config.NetworkRange = "10.9.0.0/24"
		config.DNS = "1.1.1.1,1.0.0.1"
		config.ConfigPath = "/etc/openvpn/server.conf"
	case PPTP:
		config.ListenPort = 1723
		config.NetworkRange = "10.10.0.0/24"
		config.ConfigPath = "/etc/pptpd.conf"
	case L2TP:
		config.ListenPort = 1701
		config.NetworkRange = "10.11.0.0/24"
		config.ConfigPath = "/etc/xl2tpd/xl2tpd.conf"
	}

	return config
}

// MarkProtocolEnabled marks a protocol as enabled in the database
func (vm *VPNManager) MarkProtocolEnabled(protocol Protocol, enabled bool) error {
	config, err := vm.GetOrCreateProtocolConfig(protocol)
	if err != nil {
		return err
	}

	config.Enabled = enabled
	return vm.SaveProtocolConfig(config)
}

// GetActiveConnections returns all active VPN connections from the database
func (vm *VPNManager) GetActiveConnections() ([]models.VPNConnection, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var connections []models.VPNConnection
	err := vm.db.Where("active = ?", true).
		Preload("User").
		Preload("Peer").
		Preload("Certificate").
		Order("connected_at DESC").
		Find(&connections).Error

	return connections, err
}

// GetConnectionsByProtocol returns active connections for a specific protocol
func (vm *VPNManager) GetConnectionsByProtocol(protocol Protocol) ([]models.VPNConnection, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var connections []models.VPNConnection
	err := vm.db.Where("active = ? AND protocol = ?", true, string(protocol)).
		Preload("User").
		Preload("Peer").
		Preload("Certificate").
		Order("connected_at DESC").
		Find(&connections).Error

	return connections, err
}

// WireGuard peer management methods

// CreateWireGuardPeer creates a new WireGuard peer with database persistence
func (vm *VPNManager) CreateWireGuardPeer(name, allowedIPs string) (*models.VPNPeer, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Get or create protocol config
	protocolConfig, err := vm.GetOrCreateProtocolConfig(WireGuard)
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol config: %w", err)
	}

	// Generate peer keys
	privateKeyResult, err := vm.shell.Execute("wg", "genkey")
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	privateKey := strings.TrimSpace(privateKeyResult.Stdout)

	// Generate public key from private key using pipe
	pubKeyCmd := fmt.Sprintf("echo '%s' | wg pubkey", privateKey)
	publicKeyResult, err := vm.shell.Execute("sh", "-c", pubKeyCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %w", err)
	}
	publicKey := strings.TrimSpace(publicKeyResult.Stdout)

	// Generate preshared key for added security
	presharedKeyResult, err := vm.shell.Execute("wg", "genpsk")
	if err != nil {
		return nil, fmt.Errorf("failed to generate preshared key: %w", err)
	}
	presharedKey := strings.TrimSpace(presharedKeyResult.Stdout)

	// Generate unique ID
	id := fmt.Sprintf("wg-%d", time.Now().Unix())

	// Create peer in database
	peer := &models.VPNPeer{
		ID:           id,
		ProtocolID:   protocolConfig.ID,
		Name:         name,
		PublicKey:    publicKey,
		PrivateKey:   privateKey,
		PresharedKey: presharedKey,
		AllowedIPs:   allowedIPs,
		Enabled:      true,
	}

	if err := vm.db.Create(peer).Error; err != nil {
		return nil, fmt.Errorf("failed to save peer to database: %w", err)
	}

	// Add peer to WireGuard configuration
	if err := vm.addPeerToWireGuardConfig(peer); err != nil {
		// Rollback database creation
		vm.db.Delete(peer)
		return nil, fmt.Errorf("failed to add peer to config: %w", err)
	}

	// Reload WireGuard configuration
	vm.shell.Execute("wg-quick", "down", "wg0")
	vm.shell.Execute("wg-quick", "up", "wg0")

	logger.Info("WireGuard peer created successfully",
		zap.String("peer_id", id),
		zap.String("name", name))

	return peer, nil
}

// GetWireGuardPeers returns all WireGuard peers from the database
func (vm *VPNManager) GetWireGuardPeers() ([]models.VPNPeer, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Get protocol config
	protocolConfig, err := vm.LoadProtocolConfig(WireGuard)
	if err != nil {
		return []models.VPNPeer{}, nil // Return empty list if not configured
	}

	var peers []models.VPNPeer
	err = vm.db.Where("protocol_id = ?", protocolConfig.ID).
		Order("created_at DESC").
		Find(&peers).Error

	return peers, err
}

// GetWireGuardPeer gets a specific peer by ID
func (vm *VPNManager) GetWireGuardPeer(peerID string) (*models.VPNPeer, error) {
	if vm.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var peer models.VPNPeer
	if err := vm.db.Where("id = ?", peerID).First(&peer).Error; err != nil {
		return nil, err
	}

	return &peer, nil
}

// DeleteWireGuardPeer deletes a WireGuard peer
func (vm *VPNManager) DeleteWireGuardPeer(peerID string) error {
	if vm.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Get peer from database
	peer, err := vm.GetWireGuardPeer(peerID)
	if err != nil {
		return fmt.Errorf("peer not found: %w", err)
	}

	// Remove from WireGuard configuration
	if err := vm.removePeerFromWireGuardConfig(peer.PublicKey); err != nil {
		logger.Warn("Failed to remove peer from config", zap.Error(err))
	}

	// Delete from database
	if err := vm.db.Delete(peer).Error; err != nil {
		return fmt.Errorf("failed to delete peer from database: %w", err)
	}

	// Reload WireGuard configuration
	vm.shell.Execute("wg-quick", "down", "wg0")
	vm.shell.Execute("wg-quick", "up", "wg0")

	logger.Info("WireGuard peer deleted successfully", zap.String("peer_id", peerID))
	return nil
}

// GetWireGuardPeerConfig generates client configuration for a peer
func (vm *VPNManager) GetWireGuardPeerConfig(peerID string) (string, error) {
	peer, err := vm.GetWireGuardPeer(peerID)
	if err != nil {
		return "", err
	}

	protocolConfig, err := vm.LoadProtocolConfig(WireGuard)
	if err != nil {
		return "", fmt.Errorf("protocol not configured: %w", err)
	}

	// Get server public key
	serverPublicKey := protocolConfig.PublicKey
	if serverPublicKey == "" {
		// Extract from config file
		result, _ := vm.shell.Execute("cat", "/etc/wireguard/wg0.conf")
		lines := strings.Split(result.Stdout, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PrivateKey") {
				privKey := strings.TrimSpace(strings.TrimPrefix(line, "PrivateKey = "))
				pubKeyCmd := fmt.Sprintf("echo '%s' | wg pubkey", privKey)
				pubKeyResult, _ := vm.shell.Execute("sh", "-c", pubKeyCmd)
				serverPublicKey = strings.TrimSpace(pubKeyResult.Stdout)
				break
			}
		}
	}

	// Generate client config
	serverAddress := protocolConfig.ServerAddress
	if serverAddress == "" {
		serverAddress = "YOUR_SERVER_IP"
	}

	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
DNS = %s

[Peer]
PublicKey = %s
PresharedKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25
`, peer.PrivateKey, peer.AllowedIPs, protocolConfig.DNS,
		serverPublicKey, peer.PresharedKey, serverAddress, protocolConfig.ListenPort)

	return config, nil
}

// addPeerToWireGuardConfig adds a peer to the WireGuard configuration file
func (vm *VPNManager) addPeerToWireGuardConfig(peer *models.VPNPeer) error {
	peerConfig := fmt.Sprintf(`

# Peer: %s (ID: %s)
[Peer]
PublicKey = %s
PresharedKey = %s
AllowedIPs = %s
`, peer.Name, peer.ID, peer.PublicKey, peer.PresharedKey, peer.AllowedIPs)

	// Append to config file
	cmd := fmt.Sprintf("echo '%s' >> /etc/wireguard/wg0.conf", peerConfig)
	_, err := vm.shell.Execute("sh", "-c", cmd)
	return err
}

// removePeerFromWireGuardConfig removes a peer from the WireGuard configuration file
func (vm *VPNManager) removePeerFromWireGuardConfig(publicKey string) error {
	// Read current config
	result, err := vm.shell.Execute("cat", "/etc/wireguard/wg0.conf")
	if err != nil {
		return err
	}

	// Filter out the peer section
	lines := strings.Split(result.Stdout, "\n")
	var newConfig []string
	skipPeer := false

	for _, line := range lines {
		if strings.Contains(line, "[Peer]") {
			skipPeer = false
		}
		if strings.Contains(line, "PublicKey = "+publicKey) {
			skipPeer = true
			// Remove the [Peer] line and comment before it
			if len(newConfig) >= 2 {
				newConfig = newConfig[:len(newConfig)-2]
			}
			continue
		}
		if !skipPeer {
			newConfig = append(newConfig, line)
		} else if strings.HasPrefix(line, "[") {
			skipPeer = false
			newConfig = append(newConfig, line)
		}
	}

	// Write new config
	configContent := strings.Join(newConfig, "\n")
	cmd := fmt.Sprintf("echo '%s' > /etc/wireguard/wg0.conf", configContent)
	_, err = vm.shell.Execute("sh", "-c", cmd)
	return err
}
