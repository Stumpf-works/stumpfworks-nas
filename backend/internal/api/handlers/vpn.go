package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-nas/internal/system/vpn"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/errors"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var vpnManager *vpn.VPNManager

// InitVPNManager initializes the VPN manager
func InitVPNManager(manager *vpn.VPNManager) {
	vpnManager = manager
	logger.Info("VPN manager initialized in handlers")
}

// GetVPNStatus returns the status of all VPN protocols
func GetVPNStatus(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	statuses, err := vpnManager.GetAllProtocolStatuses()
	if err != nil {
		logger.Error("Failed to get VPN status", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get VPN status", err))
		return
	}

	utils.RespondSuccess(w, statuses)
}

// GetProtocolStatus returns the status of a specific protocol
func GetProtocolStatus(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	protocolName := chi.URLParam(r, "protocol")
	if protocolName == "" {
		utils.RespondError(w, errors.BadRequest("Protocol name is required", nil))
		return
	}

	protocol := vpn.Protocol(protocolName)
	status, err := vpnManager.GetProtocolStatus(protocol)
	if err != nil {
		logger.Error("Failed to get protocol status", zap.String("protocol", protocolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get protocol status", err))
		return
	}

	utils.RespondSuccess(w, status)
}

// InstallProtocol installs a VPN protocol
func InstallProtocol(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	protocolName := chi.URLParam(r, "protocol")
	if protocolName == "" {
		utils.RespondError(w, errors.BadRequest("Protocol name is required", nil))
		return
	}

	protocol := vpn.Protocol(protocolName)
	logger.Info("Installing VPN protocol via API", zap.String("protocol", protocolName))

	// Check if already installed
	status, err := vpnManager.GetProtocolStatus(protocol)
	if err == nil && status.Installed {
		utils.RespondError(w, errors.BadRequest("Protocol already installed", nil))
		return
	}

	// Install packages
	if err := vpnManager.InstallProtocol(protocol); err != nil {
		logger.Error("Failed to install protocol", zap.String("protocol", protocolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to install protocol", err))
		return
	}

	// Initialize protocol configuration
	if err := vpnManager.InitializeProtocol(protocol); err != nil {
		logger.Error("Failed to initialize protocol", zap.String("protocol", protocolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to initialize protocol", err))
		return
	}

	logger.Info("Protocol installed and initialized successfully", zap.String("protocol", protocolName))
	utils.RespondSuccess(w, map[string]string{
		"message":  "Protocol installed and initialized successfully",
		"protocol": protocolName,
	})
}

// EnableProtocol enables and starts a VPN protocol
func EnableProtocol(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	protocolName := chi.URLParam(r, "protocol")
	if protocolName == "" {
		utils.RespondError(w, errors.BadRequest("Protocol name is required", nil))
		return
	}

	protocol := vpn.Protocol(protocolName)
	logger.Info("Enabling VPN protocol via API", zap.String("protocol", protocolName))

	// Check if installed
	status, err := vpnManager.GetProtocolStatus(protocol)
	if err != nil || !status.Installed {
		utils.RespondError(w, errors.BadRequest("Protocol not installed. Please install it first.", nil))
		return
	}

	if err := vpnManager.EnableProtocol(protocol); err != nil {
		logger.Error("Failed to enable protocol", zap.String("protocol", protocolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to enable protocol", err))
		return
	}

	logger.Info("Protocol enabled successfully", zap.String("protocol", protocolName))
	utils.RespondSuccess(w, map[string]string{
		"message":  "Protocol enabled successfully",
		"protocol": protocolName,
	})
}

// DisableProtocol stops and disables a VPN protocol
func DisableProtocol(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	protocolName := chi.URLParam(r, "protocol")
	if protocolName == "" {
		utils.RespondError(w, errors.BadRequest("Protocol name is required", nil))
		return
	}

	protocol := vpn.Protocol(protocolName)
	logger.Info("Disabling VPN protocol via API", zap.String("protocol", protocolName))

	if err := vpnManager.DisableProtocol(protocol); err != nil {
		logger.Error("Failed to disable protocol", zap.String("protocol", protocolName), zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to disable protocol", err))
		return
	}

	logger.Info("Protocol disabled successfully", zap.String("protocol", protocolName))
	utils.RespondSuccess(w, map[string]string{
		"message":  "Protocol disabled successfully",
		"protocol": protocolName,
	})
}

// WireGuard-specific handlers (for future implementation)

// CreateWireGuardPeer creates a new WireGuard peer
func CreateWireGuardPeer(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	var req struct {
		Name       string `json:"name"`
		AllowedIPs string `json:"allowed_ips"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	if req.Name == "" {
		utils.RespondError(w, errors.BadRequest("Peer name is required", nil))
		return
	}

	if req.AllowedIPs == "" {
		utils.RespondError(w, errors.BadRequest("Allowed IPs are required", nil))
		return
	}

	logger.Info("Creating WireGuard peer", zap.String("name", req.Name))

	peer, err := vpnManager.CreateWireGuardPeer(req.Name, req.AllowedIPs)
	if err != nil {
		logger.Error("Failed to create WireGuard peer", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to create peer", err))
		return
	}

	utils.RespondSuccess(w, peer)
}

// GetWireGuardPeers lists all WireGuard peers
func GetWireGuardPeers(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	peers, err := vpnManager.GetWireGuardPeers()
	if err != nil {
		logger.Error("Failed to get WireGuard peers", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get peers", err))
		return
	}

	utils.RespondSuccess(w, peers)
}

// DeleteWireGuardPeer deletes a WireGuard peer
func DeleteWireGuardPeer(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	peerID := chi.URLParam(r, "id")
	if peerID == "" {
		utils.RespondError(w, errors.BadRequest("Peer ID is required", nil))
		return
	}

	logger.Info("Deleting WireGuard peer", zap.String("peer_id", peerID))

	if err := vpnManager.DeleteWireGuardPeer(peerID); err != nil {
		logger.Error("Failed to delete WireGuard peer", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to delete peer", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"message": "Peer deleted successfully",
		"peer_id": peerID,
	})
}

// GetWireGuardPeerConfig gets the configuration file for a peer
func GetWireGuardPeerConfig(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	peerID := chi.URLParam(r, "id")
	if peerID == "" {
		utils.RespondError(w, errors.BadRequest("Peer ID is required", nil))
		return
	}

	config, err := vpnManager.GetWireGuardPeerConfig(peerID)
	if err != nil {
		logger.Error("Failed to get peer config", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get peer config", err))
		return
	}

	utils.RespondSuccess(w, map[string]string{
		"config": config,
	})
}

// GetWireGuardPeerQRCode generates a QR code for a peer configuration
func GetWireGuardPeerQRCode(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	peerID := chi.URLParam(r, "id")
	if peerID == "" {
		utils.RespondError(w, errors.BadRequest("Peer ID is required", nil))
		return
	}

	config, err := vpnManager.GetWireGuardPeerConfig(peerID)
	if err != nil {
		logger.Error("Failed to get peer config for QR code", zap.Error(err))
		utils.RespondError(w, errors.InternalServerError("Failed to get peer config", err))
		return
	}

	// Generate QR code using qrencode (installed with WireGuard)
	// Output as PNG base64
	utils.RespondSuccess(w, map[string]string{
		"qrcode": config, // TODO: Generate actual QR code using qrencode
		"config": config,
	})
}

// OpenVPN-specific handlers (for future implementation)

// GetOpenVPNCertificates lists OpenVPN certificates
func GetOpenVPNCertificates(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	// TODO: Implement certificate listing
	utils.RespondSuccess(w, []map[string]string{})
}

// CreateOpenVPNCertificate creates a new OpenVPN client certificate
func CreateOpenVPNCertificate(w http.ResponseWriter, r *http.Request) {
	if vpnManager == nil {
		utils.RespondError(w, errors.InternalServerError("VPN manager not initialized", nil))
		return
	}

	var req struct {
		CommonName string `json:"common_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, errors.BadRequest("Invalid request body", err))
		return
	}

	// TODO: Implement certificate creation
	logger.Info("Creating OpenVPN certificate", zap.String("common_name", req.CommonName))

	utils.RespondSuccess(w, map[string]string{
		"message":     "Certificate creation - to be implemented",
		"common_name": req.CommonName,
	})
}
