package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a WebSocket client
type Client struct {
	conn          *websocket.Conn
	send          chan []byte
	subscriptions map[string]bool // tracks subscribed channels
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:          conn,
		send:          make(chan []byte, 256),
		subscriptions: make(map[string]bool),
	}
}

// Read reads messages from the WebSocket connection
func (c *Client) Read() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		// Handle incoming message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Error("Failed to unmarshal WebSocket message", zap.Error(err))
			continue
		}

		// Process message based on type
		c.handleMessage(&msg)
	}
}

// Write writes messages to the WebSocket connection
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send sends a message to the client
func (c *Client) Send(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
	default:
		logger.Warn("Client send buffer full, dropping message")
	}

	return nil
}

// handleMessage handles incoming messages from the client
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case "subscribe":
		// Add channel to subscriptions
		if msg.Channel != "" {
			c.subscriptions[msg.Channel] = true
			logger.Info("Client subscribed", zap.String("channel", msg.Channel))

			// Send confirmation
			c.Send(&Message{
				Type:    "subscribed",
				Channel: msg.Channel,
			})
		}

	case "unsubscribe":
		// Remove channel from subscriptions
		if msg.Channel != "" {
			delete(c.subscriptions, msg.Channel)
			logger.Info("Client unsubscribed", zap.String("channel", msg.Channel))

			// Send confirmation
			c.Send(&Message{
				Type:    "unsubscribed",
				Channel: msg.Channel,
			})
		}

	case "ping":
		// Respond with pong
		c.Send(&Message{
			Type: "pong",
		})

	default:
		logger.Warn("Unknown message type", zap.String("type", msg.Type))
	}
}

// IsSubscribed checks if the client is subscribed to a channel
func (c *Client) IsSubscribed(channel string) bool {
	return c.subscriptions[channel]
}

// GetSubscriptions returns all channels the client is subscribed to
func (c *Client) GetSubscriptions() []string {
	channels := make([]string, 0, len(c.subscriptions))
	for channel := range c.subscriptions {
		channels = append(channels, channel)
	}
	return channels
}
