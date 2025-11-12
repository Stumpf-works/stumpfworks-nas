package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	ws "github.com/stumpfworks/nas/internal/api/websocket"
	"github.com/stumpfworks/nas/pkg/logger"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking in production
		return true
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
