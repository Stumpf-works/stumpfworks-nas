package api

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/internal/core"
)

// Handler contains all HTTP handlers
type Handler struct {
	vpnManager *core.VPNManager
}

// NewHandler creates a new handler instance
func NewHandler(vpnManager *core.VPNManager) *Handler {
	return &Handler{
		vpnManager: vpnManager,
	}
}

// General Status Handlers

// GetStatus returns the overall VPN server status
// GET /api/vpn/status
func (h *Handler) GetStatus(c *gin.Context) {
	status := h.vpnManager.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// StartServer starts all VPN servers
// POST /api/vpn/start
func (h *Handler) StartServer(c *gin.Context) {
	if err := h.vpnManager.Start(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "VPN servers started successfully",
	})
}

// StopServer stops all VPN servers
// POST /api/vpn/stop
func (h *Handler) StopServer(c *gin.Context) {
	if err := h.vpnManager.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "VPN servers stopped successfully",
	})
}

// Protocol Control Handlers

// StartProtocol starts a specific protocol
// POST /api/vpn/:protocol/start
func (h *Handler) StartProtocol(c *gin.Context) {
	protocol := c.Param("protocol")

	if err := h.vpnManager.StartProtocol(protocol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": protocol + " server started successfully",
	})
}

// StopProtocol stops a specific protocol
// POST /api/vpn/:protocol/stop
func (h *Handler) StopProtocol(c *gin.Context) {
	protocol := c.Param("protocol")

	if err := h.vpnManager.StopProtocol(protocol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": protocol + " server stopped successfully",
	})
}

// User Management Handlers

// GetUsers returns all VPN users
// GET /api/vpn/users
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.vpnManager.GetUserManager().GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// GetUser returns a specific user
// GET /api/vpn/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.vpnManager.GetUserManager().GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// CreateUser creates a new VPN user
// POST /api/vpn/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	user, err := h.vpnManager.GetUserManager().CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUserProtocols updates protocol access for a user
// PUT /api/vpn/users/:id/protocols
func (h *Handler) UpdateUserProtocols(c *gin.Context) {
	userID := c.Param("id")

	var req core.ProtocolAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := h.vpnManager.GetUserManager().UpdateProtocolAccess(userID, req.ToMap()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Protocol access updated successfully",
	})
}

// DeleteUser deletes a user
// DELETE /api/vpn/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.vpnManager.GetUserManager().DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// EnableUser enables a user account
// POST /api/vpn/users/:id/enable
func (h *Handler) EnableUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.vpnManager.GetUserManager().EnableUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User enabled successfully",
	})
}

// DisableUser disables a user account
// POST /api/vpn/users/:id/disable
func (h *Handler) DisableUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.vpnManager.GetUserManager().DisableUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User disabled successfully",
	})
}

// WireGuard Handlers

// GetWireGuardPeers returns all WireGuard peers
// GET /api/vpn/wireguard/peers
func (h *Handler) GetWireGuardPeers(c *gin.Context) {
	peers, err := h.vpnManager.GetWireGuardPeerManager().GetAllPeers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    peers,
	})
}

// GetWireGuardPeer returns a specific peer
// GET /api/vpn/wireguard/peers/:id
func (h *Handler) GetWireGuardPeer(c *gin.Context) {
	peerID := c.Param("id")

	peer, err := h.vpnManager.GetWireGuardPeerManager().GetPeer(peerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Peer not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    peer,
	})
}

// CreateWireGuardPeer creates a new WireGuard peer
// POST /api/vpn/wireguard/peers
func (h *Handler) CreateWireGuardPeer(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		AllowedIPs string  `json:"allowedIPs"`
		UserID     *string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Auto-assign IP if not provided
	if req.AllowedIPs == "" {
		nextIP, err := h.vpnManager.GetWireGuardPeerManager().GetNextAvailableIP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to allocate IP: " + err.Error(),
			})
			return
		}
		req.AllowedIPs = nextIP
	}

	peer, err := h.vpnManager.GetWireGuardPeerManager().CreatePeer(req.Name, req.AllowedIPs, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    peer,
	})
}

// DeleteWireGuardPeer deletes a peer
// DELETE /api/vpn/wireguard/peers/:id
func (h *Handler) DeleteWireGuardPeer(c *gin.Context) {
	peerID := c.Param("id")

	if err := h.vpnManager.GetWireGuardPeerManager().DeletePeer(peerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Peer deleted successfully",
	})
}

// GetWireGuardPeerConfig returns the configuration for a peer
// GET /api/vpn/wireguard/peers/:id/config
func (h *Handler) GetWireGuardPeerConfig(c *gin.Context) {
	peerID := c.Param("id")

	config, err := h.vpnManager.GetWireGuardPeerManager().GenerateConfig(peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"config": config,
		},
	})
}

// GetWireGuardPeerQRCode returns a QR code for a peer
// GET /api/vpn/wireguard/peers/:id/qrcode
func (h *Handler) GetWireGuardPeerQRCode(c *gin.Context) {
	peerID := c.Param("id")

	qrcode, err := h.vpnManager.GetWireGuardPeerManager().GenerateQRCode(peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Return as base64 encoded image
	encoded := base64.StdEncoding.EncodeToString(qrcode)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"qrcode": "data:image/png;base64," + encoded,
		},
	})
}

// GetWireGuardStats returns statistics for WireGuard peers
// GET /api/vpn/wireguard/stats
func (h *Handler) GetWireGuardStats(c *gin.Context) {
	// Get peer stats from server
	// This is a placeholder - in production, implement proper stats collection
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Stats endpoint - to be implemented",
		},
	})
}

// SearchUsers searches for users
// GET /api/vpn/users/search?q=query
func (h *Handler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Query parameter 'q' is required",
		})
		return
	}

	users, err := h.vpnManager.GetUserManager().SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}
