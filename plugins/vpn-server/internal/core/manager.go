package core

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/config"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/internal/wireguard"
)

// VPNManager is the central manager for all VPN protocols
type VPNManager struct {
	config      *config.Config
	db          *gorm.DB
	wireguard   *wireguard.Server
	wgPeerMgr   *wireguard.PeerManager
	userManager *UserManager
	// openvpn     *openvpn.Server   // TODO: Phase 2
	// pptp        *pptp.Server      // TODO: Phase 3
	// l2tp        *l2tp.Server      // TODO: Phase 3
	running bool
	mu      sync.RWMutex
}

// NewVPNManager creates a new VPN manager instance
func NewVPNManager(cfg *config.Config, db *gorm.DB) (*VPNManager, error) {
	mgr := &VPNManager{
		config:      cfg,
		db:          db,
		userManager: NewUserManager(db),
		running:     false,
	}

	// Initialize WireGuard if enabled
	if cfg.WireGuard.Enabled {
		wg, err := wireguard.NewServer(&cfg.WireGuard, db)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize WireGuard: %w", err)
		}
		mgr.wireguard = wg
		mgr.wgPeerMgr = wireguard.NewPeerManager(db, wg)
	}

	// TODO: Initialize other protocols in future phases
	// if cfg.OpenVPN.Enabled {
	//     ovpn, err := openvpn.NewServer(&cfg.OpenVPN, db)
	//     if err != nil {
	//         return nil, fmt.Errorf("failed to initialize OpenVPN: %w", err)
	//     }
	//     mgr.openvpn = ovpn
	// }

	return mgr, nil
}

// Start starts all enabled VPN servers
func (m *VPNManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("VPN manager already running")
	}

	// Start WireGuard
	if m.wireguard != nil {
		if err := m.wireguard.Start(); err != nil {
			return fmt.Errorf("failed to start WireGuard: %w", err)
		}
	}

	// TODO: Start other protocols
	// if m.openvpn != nil {
	//     if err := m.openvpn.Start(); err != nil {
	//         return fmt.Errorf("failed to start OpenVPN: %w", err)
	//     }
	// }

	m.running = true
	return nil
}

// Stop stops all VPN servers
func (m *VPNManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("VPN manager not running")
	}

	// Stop WireGuard
	if m.wireguard != nil {
		if err := m.wireguard.Stop(); err != nil {
			return fmt.Errorf("failed to stop WireGuard: %w", err)
		}
	}

	// TODO: Stop other protocols

	m.running = false
	return nil
}

// IsRunning returns whether the VPN manager is running
func (m *VPNManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// StartProtocol starts a specific protocol server
func (m *VPNManager) StartProtocol(protocol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch protocol {
	case "wireguard":
		if m.wireguard == nil {
			return fmt.Errorf("WireGuard not initialized")
		}
		return m.wireguard.Start()
	// TODO: Add other protocols
	default:
		return fmt.Errorf("unknown protocol: %s", protocol)
	}
}

// StopProtocol stops a specific protocol server
func (m *VPNManager) StopProtocol(protocol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch protocol {
	case "wireguard":
		if m.wireguard == nil {
			return fmt.Errorf("WireGuard not initialized")
		}
		return m.wireguard.Stop()
	// TODO: Add other protocols
	default:
		return fmt.Errorf("unknown protocol: %s", protocol)
	}
}

// GetStatus returns the status of all VPN servers
func (m *VPNManager) GetStatus() *VPNStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := &VPNStatus{
		Running: m.running,
	}

	// WireGuard status
	if m.wireguard != nil {
		status.WireGuard = m.wireguard.GetStatus()
	}

	// TODO: Add other protocol statuses

	// Calculate statistics
	status.Statistics = m.calculateStatistics()

	return status
}

// GetUserManager returns the user manager
func (m *VPNManager) GetUserManager() *UserManager {
	return m.userManager
}

// GetWireGuardPeerManager returns the WireGuard peer manager
func (m *VPNManager) GetWireGuardPeerManager() *wireguard.PeerManager {
	return m.wgPeerMgr
}

// Close cleans up all resources
func (m *VPNManager) Close() error {
	if m.running {
		if err := m.Stop(); err != nil {
			return err
		}
	}

	if m.wireguard != nil {
		if err := m.wireguard.Close(); err != nil {
			return err
		}
	}

	// TODO: Close other protocols

	return nil
}

// Helper methods

func (m *VPNManager) calculateStatistics() *Statistics {
	stats := &Statistics{}

	// Count active protocols
	if m.wireguard != nil && m.wireguard.IsRunning() {
		stats.ActiveProtocols++
	}

	// Get WireGuard stats
	if m.wireguard != nil {
		wgStatus := m.wireguard.GetStatus()
		stats.TotalConnections += wgStatus.ConnectedPeers
		stats.TotalBytesIn += wgStatus.BytesReceived
		stats.TotalBytesOut += wgStatus.BytesSent
		stats.ConnectionsByProtocol["wireguard"] = wgStatus.ConnectedPeers
	}

	// TODO: Add stats from other protocols

	return stats
}

// Types for API responses

// VPNStatus represents the overall VPN server status
type VPNStatus struct {
	Running    bool                    `json:"running"`
	WireGuard  *wireguard.ServerStatus `json:"wireguard,omitempty"`
	// OpenVPN    *openvpn.ServerStatus   `json:"openvpn,omitempty"`   // TODO
	// PPTP       *pptp.ServerStatus      `json:"pptp,omitempty"`      // TODO
	// L2TP       *l2tp.ServerStatus      `json:"l2tp,omitempty"`      // TODO
	Statistics *Statistics             `json:"statistics"`
}

// Statistics represents aggregated VPN statistics
type Statistics struct {
	TotalConnections       int            `json:"totalConnections"`
	ActiveProtocols        int            `json:"activeProtocols"`
	TotalBytesIn           int64          `json:"totalBytesIn"`
	TotalBytesOut          int64          `json:"totalBytesOut"`
	ConnectionsByProtocol  map[string]int `json:"connectionsByProtocol"`
}

func init() {
	// Initialize maps
	_ = &Statistics{
		ConnectionsByProtocol: make(map[string]int),
	}
}
