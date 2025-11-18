package ami

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Client represents an Asterisk Manager Interface client
type Client struct {
	host     string
	port     int
	username string
	secret   string

	conn      net.Conn
	connected bool
	mu        sync.RWMutex

	eventListeners []chan Event
}

// Event represents an AMI event
type Event struct {
	Name   string
	Fields map[string]string
}

// Response represents an AMI response
type Response struct {
	Success bool
	Message string
	Fields  map[string]string
}

// NewClient creates a new AMI client
func NewClient(host string, port int, username, secret string) *Client {
	return &Client{
		host:           host,
		port:           port,
		username:       username,
		secret:         secret,
		eventListeners: make([]chan Event, 0),
	}
}

// Connect establishes connection to Asterisk AMI
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	address := fmt.Sprintf("%s:%d", c.host, c.port)
	log.Info().Str("address", address).Msg("Connecting to Asterisk AMI")

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to AMI: %w", err)
	}

	c.conn = conn

	// Read welcome banner
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to read banner: %w", err)
	}

	log.Debug().Str("banner", strings.TrimSpace(banner)).Msg("Received AMI banner")

	// Login
	if err := c.login(); err != nil {
		conn.Close()
		return fmt.Errorf("login failed: %w", err)
	}

	c.connected = true
	log.Info().Msg("Successfully connected to Asterisk AMI")

	// Start event listener
	go c.eventLoop()

	return nil
}

// login performs AMI login
func (c *Client) login() error {
	loginCmd := fmt.Sprintf(
		"Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n",
		c.username,
		c.secret,
	)

	_, err := c.conn.Write([]byte(loginCmd))
	if err != nil {
		return fmt.Errorf("failed to send login: %w", err)
	}

	// Read response
	reader := bufio.NewReader(c.conn)
	resp, err := c.readResponse(reader)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("login failed: %s", resp.Message)
	}

	log.Info().Msg("AMI login successful")
	return nil
}

// Disconnect closes the AMI connection
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	// Send Logoff action
	logoffCmd := "Action: Logoff\r\n\r\n"
	c.conn.Write([]byte(logoffCmd))

	c.conn.Close()
	c.connected = false

	log.Info().Msg("Disconnected from Asterisk AMI")
	return nil
}

// IsConnected returns connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// SendAction sends an action to Asterisk
func (c *Client) SendAction(action string, params map[string]string) (*Response, error) {
	c.mu.RLock()
	if !c.connected {
		c.mu.RUnlock()
		return nil, fmt.Errorf("not connected to AMI")
	}
	c.mu.RUnlock()

	// Build action string
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Action: %s\r\n", action))

	for key, value := range params {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	sb.WriteString("\r\n")

	// Send action
	_, err := c.conn.Write([]byte(sb.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to send action: %w", err)
	}

	// Read response
	reader := bufio.NewReader(c.conn)
	return c.readResponse(reader)
}

// readResponse reads and parses AMI response
func (c *Client) readResponse(reader *bufio.Reader) (*Response, error) {
	fields := make(map[string]string)
	var message string
	success := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		line = strings.TrimSpace(line)

		// Empty line indicates end of response
		if line == "" {
			break
		}

		// Parse key-value pair
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		fields[key] = value

		// Check response status
		if key == "Response" {
			success = value == "Success"
		}

		if key == "Message" {
			message = value
		}
	}

	return &Response{
		Success: success,
		Message: message,
		Fields:  fields,
	}, nil
}

// eventLoop continuously reads events from AMI
func (c *Client) eventLoop() {
	reader := bufio.NewReader(c.conn)

	for c.IsConnected() {
		event, err := c.readEvent(reader)
		if err != nil {
			if c.IsConnected() {
				log.Error().Err(err).Msg("Error reading event")
			}
			return
		}

		if event != nil {
			c.broadcastEvent(*event)
		}
	}
}

// readEvent reads an event from AMI
func (c *Client) readEvent(reader *bufio.Reader) (*Event, error) {
	fields := make(map[string]string)
	var eventName string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)

		// Empty line indicates end of event
		if line == "" {
			if eventName != "" {
				return &Event{
					Name:   eventName,
					Fields: fields,
				}, nil
			}
			continue
		}

		// Parse key-value pair
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		fields[key] = value

		if key == "Event" {
			eventName = value
		}
	}
}

// SubscribeEvents subscribes to AMI events
func (c *Client) SubscribeEvents() chan Event {
	eventChan := make(chan Event, 100)
	c.mu.Lock()
	c.eventListeners = append(c.eventListeners, eventChan)
	c.mu.Unlock()
	return eventChan
}

// broadcastEvent sends event to all listeners
func (c *Client) broadcastEvent(event Event) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, listener := range c.eventListeners {
		select {
		case listener <- event:
		default:
			// Channel full, skip
			log.Warn().Str("event", event.Name).Msg("Event listener channel full")
		}
	}
}

// CoreShowVersion returns Asterisk version
func (c *Client) CoreShowVersion() (string, error) {
	resp, err := c.SendAction("Command", map[string]string{
		"Command": "core show version",
	})
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", fmt.Errorf("command failed: %s", resp.Message)
	}

	return resp.Fields["Output"], nil
}

// SIPPeers returns list of SIP peers
func (c *Client) SIPPeers() ([]map[string]string, error) {
	resp, err := c.SendAction("SIPpeers", nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("command failed: %s", resp.Message)
	}

	// Parse peers (simplified - real implementation would read multiple events)
	peers := make([]map[string]string, 0)
	return peers, nil
}

// Originate initiates a call
func (c *Client) Originate(channel, extension, context string, timeout int) error {
	_, err := c.SendAction("Originate", map[string]string{
		"Channel":  channel,
		"Exten":    extension,
		"Context":  context,
		"Priority": "1",
		"Timeout":  fmt.Sprintf("%d", timeout*1000),
	})
	return err
}

// Hangup hangs up a channel
func (c *Client) Hangup(channel string) error {
	_, err := c.SendAction("Hangup", map[string]string{
		"Channel": channel,
	})
	return err
}

// Reload reloads Asterisk configuration
func (c *Client) Reload(module string) error {
	params := make(map[string]string)
	if module != "" {
		params["Module"] = module
	}

	_, err := c.SendAction("Reload", params)
	return err
}
