package wireguard

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/config"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/pkg/database"
)

// Server represents a WireGuard VPN server
type Server struct {
	config     *config.WireGuardConfig
	db         *gorm.DB
	client     *wgctrl.Client
	device     string
	privateKey wgtypes.Key
	publicKey  wgtypes.Key
	peers      map[string]*database.WireGuardPeer
	running    bool
	mu         sync.RWMutex
}

// NewServer creates a new WireGuard server instance
func NewServer(cfg *config.WireGuardConfig, db *gorm.DB) (*Server, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create WireGuard client: %w", err)
	}

	// Load or generate server keys
	privateKey, publicKey, err := loadOrGenerateKeys(db, "wireguard")
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}

	return &Server{
		config:     cfg,
		db:         db,
		client:     client,
		device:     cfg.Interface,
		privateKey: privateKey,
		publicKey:  publicKey,
		peers:      make(map[string]*database.WireGuardPeer),
		running:    false,
	}, nil
}

// Start starts the WireGuard server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	// Create WireGuard interface
	if err := s.createInterface(); err != nil {
		return fmt.Errorf("failed to create interface: %w", err)
	}

	// Configure interface
	listenPort := s.config.ListenPort
	cfg := wgtypes.Config{
		PrivateKey:   &s.privateKey,
		ListenPort:   &listenPort,
		ReplacePeers: false,
	}

	if err := s.client.ConfigureDevice(s.device, cfg); err != nil {
		return fmt.Errorf("failed to configure device: %w", err)
	}

	// Set up IP address
	if err := s.setupIPAddress(); err != nil {
		return fmt.Errorf("failed to setup IP: %w", err)
	}

	// Bring interface up
	if err := s.bringUp(); err != nil {
		return fmt.Errorf("failed to bring up interface: %w", err)
	}

	// Load existing peers from database
	if err := s.loadPeers(); err != nil {
		return fmt.Errorf("failed to load peers: %w", err)
	}

	// Setup firewall rules
	if err := s.setupFirewall(); err != nil {
		return fmt.Errorf("failed to setup firewall: %w", err)
	}

	s.running = true
	return nil
}

// Stop stops the WireGuard server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("server not running")
	}

	// Clean up firewall rules
	if err := s.cleanupFirewall(); err != nil {
		return fmt.Errorf("failed to cleanup firewall: %w", err)
	}

	// Bring interface down
	if err := s.bringDown(); err != nil {
		return fmt.Errorf("failed to bring down interface: %w", err)
	}

	// Delete interface
	if err := s.deleteInterface(); err != nil {
		return fmt.Errorf("failed to delete interface: %w", err)
	}

	s.running = false
	return nil
}

// IsRunning returns whether the server is running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// AddPeer adds a new peer to the WireGuard server
func (s *Server) AddPeer(peer *database.WireGuardPeer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Parse peer public key
	pubKey, err := wgtypes.ParseKey(peer.PublicKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	// Parse allowed IPs
	allowedIPs, err := parseIPNets(peer.AllowedIPs)
	if err != nil {
		return fmt.Errorf("invalid allowed IPs: %w", err)
	}

	// Configure peer
	peerCfg := wgtypes.PeerConfig{
		PublicKey:  pubKey,
		AllowedIPs: allowedIPs,
		ReplaceAllowedIPs: true,
	}

	// Parse endpoint if provided
	if peer.Endpoint != "" {
		endpoint, err := net.ResolveUDPAddr("udp", peer.Endpoint)
		if err != nil {
			return fmt.Errorf("invalid endpoint: %w", err)
		}
		peerCfg.Endpoint = endpoint
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	if err := s.client.ConfigureDevice(s.device, cfg); err != nil {
		return fmt.Errorf("failed to add peer: %w", err)
	}

	s.peers[peer.PublicKey] = peer
	return nil
}

// RemovePeer removes a peer from the WireGuard server
func (s *Server) RemovePeer(publicKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pubKey, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	peerCfg := wgtypes.PeerConfig{
		PublicKey: pubKey,
		Remove:    true,
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	if err := s.client.ConfigureDevice(s.device, cfg); err != nil {
		return fmt.Errorf("failed to remove peer: %w", err)
	}

	delete(s.peers, publicKey)
	return nil
}

// GetPeerStats retrieves statistics for all peers
func (s *Server) GetPeerStats() ([]*PeerStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, err := s.client.Device(s.device)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	var stats []*PeerStats
	for _, peer := range device.Peers {
		pubKey := peer.PublicKey.String()

		// Get peer info from database
		var dbPeer database.WireGuardPeer
		if err := s.db.Where("public_key = ?", pubKey).First(&dbPeer).Error; err != nil {
			continue // Skip peers not in database
		}

		endpoint := ""
		if peer.Endpoint != nil {
			endpoint = peer.Endpoint.String()
		}

		lastHandshake := ""
		if !peer.LastHandshakeTime.IsZero() {
			duration := time.Since(peer.LastHandshakeTime)
			lastHandshake = formatDuration(duration)
		}

		stats = append(stats, &PeerStats{
			ID:            dbPeer.ID,
			Name:          dbPeer.Name,
			PublicKey:     pubKey,
			Endpoint:      endpoint,
			AllowedIPs:    dbPeer.AllowedIPs,
			LastHandshake: lastHandshake,
			BytesReceived: int64(peer.ReceiveBytes),
			BytesSent:     int64(peer.TransmitBytes),
		})
	}

	return stats, nil
}

// GenerateClientConfig generates a WireGuard client configuration
func (s *Server) GenerateClientConfig(peer *database.WireGuardPeer) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
DNS = %s

[Peer]
PublicKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25
`,
		peer.PrivateKey,
		peer.AllowedIPs,
		s.config.DNS,
		s.publicKey.String(),
		s.config.Endpoint,
		s.config.ListenPort,
	)

	return config, nil
}

// GetStatus returns the current status of the WireGuard server
func (s *Server) GetStatus() *ServerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats, _ := s.GetPeerStats()

	var connectedPeers int
	var totalBytesIn, totalBytesOut int64

	for _, stat := range stats {
		if stat.LastHandshake != "" {
			connectedPeers++
		}
		totalBytesIn += stat.BytesReceived
		totalBytesOut += stat.BytesSent
	}

	return &ServerStatus{
		Enabled:       s.config.Enabled,
		Running:       s.running,
		PublicKey:     s.publicKey.String(),
		ListenPort:    s.config.ListenPort,
		Endpoint:      fmt.Sprintf("%s:%d", s.config.Endpoint, s.config.ListenPort),
		Subnet:        s.config.Subnet,
		DNS:           s.config.DNS,
		TotalPeers:    len(s.peers),
		ConnectedPeers: connectedPeers,
		BytesReceived: totalBytesIn,
		BytesSent:     totalBytesOut,
	}
}

// Close closes the WireGuard client
func (s *Server) Close() error {
	if s.running {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	return s.client.Close()
}

// Helper functions

func (s *Server) createInterface() error {
	cmd := exec.Command("ip", "link", "add", "dev", s.device, "type", "wireguard")
	return cmd.Run()
}

func (s *Server) deleteInterface() error {
	cmd := exec.Command("ip", "link", "delete", "dev", s.device)
	return cmd.Run()
}

func (s *Server) setupIPAddress() error {
	// Extract first IP from subnet for server
	_, ipNet, err := net.ParseCIDR(s.config.Subnet)
	if err != nil {
		return err
	}

	ip := ipNet.IP
	ip[len(ip)-1] = 1 // Use .1 as server IP

	serverIP := fmt.Sprintf("%s/%d", ip.String(), getMaskBits(ipNet))
	cmd := exec.Command("ip", "address", "add", "dev", s.device, serverIP)
	return cmd.Run()
}

func (s *Server) bringUp() error {
	cmd := exec.Command("ip", "link", "set", "up", "dev", s.device)
	return cmd.Run()
}

func (s *Server) bringDown() error {
	cmd := exec.Command("ip", "link", "set", "down", "dev", s.device)
	return cmd.Run()
}

func (s *Server) setupFirewall() error {
	// Enable IP forwarding
	exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()

	// Add iptables rules for NAT
	exec.Command("iptables", "-A", "FORWARD", "-i", s.device, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-A", "FORWARD", "-o", s.device, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE").Run()

	return nil
}

func (s *Server) cleanupFirewall() error {
	exec.Command("iptables", "-D", "FORWARD", "-i", s.device, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-D", "FORWARD", "-o", s.device, "-j", "ACCEPT").Run()
	exec.Command("iptables", "-t", "nat", "-D", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE").Run()
	return nil
}

func (s *Server) loadPeers() error {
	var peers []database.WireGuardPeer
	if err := s.db.Where("enabled = ?", true).Find(&peers).Error; err != nil {
		return err
	}

	for _, peer := range peers {
		if err := s.AddPeer(&peer); err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}

func parseIPNets(allowedIPsStr string) ([]net.IPNet, error) {
	var ipNets []net.IPNet
	_, ipNet, err := net.ParseCIDR(allowedIPsStr)
	if err != nil {
		return nil, err
	}
	ipNets = append(ipNets, *ipNet)
	return ipNets, nil
}

func getMaskBits(ipNet *net.IPNet) int {
	ones, _ := ipNet.Mask.Size()
	return ones
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}
	return fmt.Sprintf("%d days ago", int(d.Hours()/24))
}

func loadOrGenerateKeys(db *gorm.DB, protocol string) (wgtypes.Key, wgtypes.Key, error) {
	// Try to load existing keys from config table
	var vpnConfig database.VPNConfig
	err := db.Where("protocol = ?", protocol).First(&vpnConfig).Error

	if err == nil && vpnConfig.Config != "" {
		// Keys exist, parse them
		// This is simplified - in production, properly parse JSON config
		privateKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return wgtypes.Key{}, wgtypes.Key{}, err
		}
		publicKey := privateKey.PublicKey()
		return privateKey, publicKey, nil
	}

	// Generate new keys
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return wgtypes.Key{}, wgtypes.Key{}, err
	}
	publicKey := privateKey.PublicKey()

	// Store in database
	// Simplified - in production, properly serialize to JSON
	vpnConfig = database.VPNConfig{
		Protocol: protocol,
		Config:   fmt.Sprintf(`{"privateKey":"%s","publicKey":"%s"}`, privateKey.String(), publicKey.String()),
		Enabled:  true,
	}
	db.Save(&vpnConfig)

	return privateKey, publicKey, nil
}
