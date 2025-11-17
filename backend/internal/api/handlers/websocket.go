// Revision: 2025-11-17 | Author: Claude | Version: 1.1.2
package handlers

import (
	"net/http"
	"strings"

	ws "github.com/Stumpf-works/stumpfworks-nas/internal/api/websocket"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// createUpgrader creates a WebSocket upgrader with origin checking based on config
func createUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// No origin header - allow (e.g., from non-browser clients, curl, etc.)
				return true
			}

			// Get allowed origins from config
			cfg := config.GlobalConfig
			if cfg == nil {
				logger.Error("Config not initialized for WebSocket origin check")
				return false
			}

			// Check against configured allowed origins
			for _, allowed := range cfg.Server.AllowedOrigins {
				if origin == allowed {
					return true
				}
				// Also check with trailing slash variants
				if origin+"/" == allowed || origin == allowed+"/" {
					return true
				}
			}

			// In development, also allow same-origin requests
			if cfg.IsDevelopment() {
				host := r.Host
				if strings.Contains(origin, host) {
					logger.Debug("WebSocket same-origin allowed in development",
						zap.String("origin", origin))
					return true
				}
			}

			// Deny and log
			logger.Warn("WebSocket connection from unauthorized origin denied",
				zap.String("origin", origin),
				zap.Strings("allowed_origins", cfg.Server.AllowedOrigins))
			return false
		},
	}
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := createUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade WebSocket connection", zap.Error(err))
		return
	}

	client := ws.NewClient(conn)
	go client.Read()
	go client.Write()

	logger.Info("WebSocket client connected",
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("origin", r.Header.Get("Origin")))
}
