package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	ws "github.com/Stumpf-works/stumpfworks-nas/internal/api/websocket"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// No origin header - allow (e.g., from non-browser clients)
			return true
		}

		// Allow localhost and 127.0.0.1 for development
		if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
			return true
		}

		// Allow same-origin requests (origin matches host)
		host := r.Host
		if strings.Contains(origin, host) {
			return true
		}

		// In production, you might want to check against a whitelist
		// For now, log and deny unknown origins
		logger.Warn("WebSocket connection from unknown origin denied",
			zap.String("origin", origin),
			zap.String("host", host))
		return false
	},
}

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	client := ws.NewClient(conn)
	go client.Read()
	go client.Write()

	logger.Info("WebSocket client connected", zap.String("remote_addr", r.RemoteAddr))
}
