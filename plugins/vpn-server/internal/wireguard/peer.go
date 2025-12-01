package wireguard

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/pkg/database"
)

// PeerManager handles WireGuard peer operations
type PeerManager struct {
	db     *gorm.DB
	server *Server
}

// NewPeerManager creates a new peer manager
func NewPeerManager(db *gorm.DB, server *Server) *PeerManager {
	return &PeerManager{
		db:     db,
		server: server,
	}
}

// CreatePeer creates a new WireGuard peer
func (pm *PeerManager) CreatePeer(name, allowedIPs string, userID *string) (*database.WireGuardPeer, error) {
	// Generate key pair
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	publicKey := privateKey.PublicKey()

	// Create peer in database
	peer := &database.WireGuardPeer{
		ID:         uuid.New().String(),
		UserID:     userID,
		Name:       name,
		PublicKey:  publicKey.String(),
		PrivateKey: privateKey.String(),
		AllowedIPs: allowedIPs,
		Enabled:    true,
	}

	if err := pm.db.Create(peer).Error; err != nil {
		return nil, fmt.Errorf("failed to create peer in database: %w", err)
	}

	// Add peer to WireGuard server if running
	if pm.server.IsRunning() {
		if err := pm.server.AddPeer(peer); err != nil {
			// Rollback database creation
			pm.db.Delete(peer)
			return nil, fmt.Errorf("failed to add peer to server: %w", err)
		}
	}

	return peer, nil
}

// GetPeer retrieves a peer by ID
func (pm *PeerManager) GetPeer(id string) (*database.WireGuardPeer, error) {
	var peer database.WireGuardPeer
	if err := pm.db.Where("id = ?", id).First(&peer).Error; err != nil {
		return nil, err
	}
	return &peer, nil
}

// GetPeerByPublicKey retrieves a peer by public key
func (pm *PeerManager) GetPeerByPublicKey(publicKey string) (*database.WireGuardPeer, error) {
	var peer database.WireGuardPeer
	if err := pm.db.Where("public_key = ?", publicKey).First(&peer).Error; err != nil {
		return nil, err
	}
	return &peer, nil
}

// GetAllPeers retrieves all peers
func (pm *PeerManager) GetAllPeers() ([]database.WireGuardPeer, error) {
	var peers []database.WireGuardPeer
	if err := pm.db.Find(&peers).Error; err != nil {
		return nil, err
	}
	return peers, nil
}

// GetPeersByUser retrieves all peers for a specific user
func (pm *PeerManager) GetPeersByUser(userID string) ([]database.WireGuardPeer, error) {
	var peers []database.WireGuardPeer
	if err := pm.db.Where("user_id = ?", userID).Find(&peers).Error; err != nil {
		return nil, err
	}
	return peers, nil
}

// UpdatePeer updates a peer
func (pm *PeerManager) UpdatePeer(peer *database.WireGuardPeer) error {
	return pm.db.Save(peer).Error
}

// DeletePeer deletes a peer
func (pm *PeerManager) DeletePeer(id string) error {
	peer, err := pm.GetPeer(id)
	if err != nil {
		return err
	}

	// Remove from WireGuard server if running
	if pm.server.IsRunning() {
		if err := pm.server.RemovePeer(peer.PublicKey); err != nil {
			return fmt.Errorf("failed to remove peer from server: %w", err)
		}
	}

	// Delete from database
	if err := pm.db.Delete(peer).Error; err != nil {
		return fmt.Errorf("failed to delete peer from database: %w", err)
	}

	return nil
}

// EnablePeer enables a peer
func (pm *PeerManager) EnablePeer(id string) error {
	peer, err := pm.GetPeer(id)
	if err != nil {
		return err
	}

	peer.Enabled = true
	if err := pm.db.Save(peer).Error; err != nil {
		return err
	}

	// Add to server if running
	if pm.server.IsRunning() {
		return pm.server.AddPeer(peer)
	}

	return nil
}

// DisablePeer disables a peer
func (pm *PeerManager) DisablePeer(id string) error {
	peer, err := pm.GetPeer(id)
	if err != nil {
		return err
	}

	peer.Enabled = false
	if err := pm.db.Save(peer).Error; err != nil {
		return err
	}

	// Remove from server if running
	if pm.server.IsRunning() {
		return pm.server.RemovePeer(peer.PublicKey)
	}

	return nil
}

// GenerateConfig generates a WireGuard configuration file for a peer
func (pm *PeerManager) GenerateConfig(id string) (string, error) {
	peer, err := pm.GetPeer(id)
	if err != nil {
		return "", err
	}

	return pm.server.GenerateClientConfig(peer)
}

// GenerateQRCode generates a QR code for a peer configuration
func (pm *PeerManager) GenerateQRCode(id string) ([]byte, error) {
	config, err := pm.GenerateConfig(id)
	if err != nil {
		return nil, err
	}

	// Generate QR code
	qr, err := qrcode.Encode(config, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qr, nil
}

// GetNextAvailableIP returns the next available IP address in the subnet
func (pm *PeerManager) GetNextAvailableIP() (string, error) {
	// Get all existing peers
	peers, err := pm.GetAllPeers()
	if err != nil {
		return "", err
	}

	// Parse subnet from server config
	subnet := pm.server.config.Subnet

	// Simple implementation - in production, use proper IP allocation
	// For now, just use .2, .3, .4, etc.
	nextIP := len(peers) + 2
	return fmt.Sprintf("10.8.0.%d/32", nextIP), nil
}

// Types for API responses

// PeerStats represents peer statistics
type PeerStats struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	PublicKey     string `json:"publicKey"`
	Endpoint      string `json:"endpoint,omitempty"`
	AllowedIPs    string `json:"allowedIPs"`
	LastHandshake string `json:"lastHandshake,omitempty"`
	BytesReceived int64  `json:"bytesReceived"`
	BytesSent     int64  `json:"bytesSent"`
}

// ServerStatus represents server status
type ServerStatus struct {
	Enabled        bool   `json:"enabled"`
	Running        bool   `json:"running"`
	PublicKey      string `json:"publicKey"`
	ListenPort     int    `json:"listenPort"`
	Endpoint       string `json:"endpoint"`
	Subnet         string `json:"subnet"`
	DNS            string `json:"dns"`
	TotalPeers     int    `json:"totalPeers"`
	ConnectedPeers int    `json:"connectedPeers"`
	BytesReceived  int64  `json:"bytesReceived"`
	BytesSent      int64  `json:"bytesSent"`
}

// PeerConfig represents a peer configuration for export
type PeerConfig struct {
	Name       string `json:"name"`
	Config     string `json:"config"`
	QRCode     string `json:"qrCode,omitempty"` // Base64 encoded
}
