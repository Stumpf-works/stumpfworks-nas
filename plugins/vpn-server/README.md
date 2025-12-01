# VPN Server Plugin for Stumpfworks NAS

A comprehensive multi-protocol VPN server application supporting WireGuard, OpenVPN, PPTP, and L2TP/IPsec with a stunning Stumpfworks-designed interface.

## 🎯 Overview

This plugin provides a unified VPN management interface that surpasses commercial NAS solutions like Synology, offering:

- **Multi-Protocol Support**: WireGuard, OpenVPN, PPTP, L2TP/IPsec
- **Unified User Management**: Cross-protocol user permissions matrix
- **Beautiful UI**: Stumpfworks design with gradients, animations, and glassmorphism
- **Real-time Monitoring**: Live connection stats and traffic analytics
- **Enterprise Features**: LDAP/AD integration, comprehensive logging

## 🏗️ Architecture

### Frontend Structure
```
frontend/src/apps/VPNServer/
├── VPNServer.tsx              # Main application entry point
├── components/
│   ├── Dashboard.tsx          # Protocol overview with stats
│   ├── ConnectionList.tsx     # User permissions matrix
│   ├── GeneralSettings.tsx    # Global VPN settings
│   └── protocols/
│       ├── WireGuardPanel.tsx # WireGuard management
│       ├── OpenVPNPanel.tsx   # OpenVPN management
│       ├── PPTPPanel.tsx      # PPTP management
│       └── L2TPPanel.tsx      # L2TP/IPsec management
└── types/
    └── vpn.ts                 # TypeScript definitions
```

### Backend Structure
```
plugins/vpn-server/
├── cmd/vpn-server/
│   └── main.go                # Plugin entry point
├── internal/
│   ├── core/
│   │   ├── manager.go         # Core VPN manager
│   │   ├── user.go            # User management
│   │   └── stats.go           # Statistics aggregation
│   ├── wireguard/
│   │   ├── server.go          # WireGuard server
│   │   ├── peer.go            # Peer management
│   │   └── config.go          # Config generation
│   ├── openvpn/
│   │   ├── server.go          # OpenVPN server
│   │   ├── certificate.go     # PKI management
│   │   └── config.go          # Config generation
│   ├── pptp/
│   │   ├── server.go          # PPTP server
│   │   └── config.go          # Config management
│   ├── l2tp/
│   │   ├── server.go          # L2TP server
│   │   ├── ipsec.go           # IPsec configuration
│   │   └── config.go          # Config management
│   └── api/
│       ├── handlers.go        # HTTP handlers
│       ├── routes.go          # Route definitions
│       └── middleware.go      # Auth & logging
├── pkg/
│   ├── database/
│   │   └── models.go          # Database models
│   └── utils/
│       ├── network.go         # Network utilities
│       └── crypto.go          # Cryptographic helpers
└── go.mod
```

## 🚀 Implementation Phases

### Phase 1: Foundation (Week 1-2)
**Goal**: Core infrastructure and WireGuard implementation

- [x] Frontend UI components
- [ ] Backend Go module setup
- [ ] Database schema for users and connections
- [ ] WireGuard server implementation
- [ ] Basic API endpoints
- [ ] User management system

### Phase 2: Protocol Expansion (Week 3-4)
**Goal**: Add OpenVPN support

- [ ] OpenVPN server implementation
- [ ] PKI/Certificate management
- [ ] Certificate generation and revocation
- [ ] OpenVPN client config generation
- [ ] Integration with user management

### Phase 3: Legacy Protocols (Week 5-6)
**Goal**: Add PPTP and L2TP/IPsec

- [ ] PPTP server implementation
- [ ] L2TP/IPsec server implementation
- [ ] IPsec PSK management
- [ ] Security warnings for PPTP
- [ ] Client setup documentation

### Phase 4: Advanced Features (Week 7-8)
**Goal**: Enterprise features and polish

- [ ] LDAP/Active Directory integration
- [ ] Advanced logging and monitoring
- [ ] Traffic analytics and graphs
- [ ] Email notifications
- [ ] Automatic client config distribution
- [ ] QR code generation for mobile clients
- [ ] Backup/restore configuration

## 📦 Dependencies

### Backend (Go)
```go
// go.mod
module github.com/stumpfworks/nas-plugins/vpn-server

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1              // HTTP framework
    github.com/spf13/viper v1.17.0               // Configuration
    golang.zx2c4.com/wireguard/wgctrl v0.0.0    // WireGuard control
    github.com/OpenVPN/openvpn3 v3.8.2           // OpenVPN
    github.com/go-sql-driver/mysql v1.7.1        // Database
    github.com/golang-jwt/jwt/v5 v5.1.0          // JWT auth
    github.com/skip2/go-qrcode v0.0.0            // QR codes
    gopkg.in/gomail.v2 v2.0.0                    // Email
)
```

### Frontend (React)
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "framer-motion": "^10.16.4",
    "lucide-react": "^0.292.0",
    "tailwindcss": "^3.3.5"
  }
}
```

### System Requirements
- **WireGuard**: Kernel module or wireguard-go
- **OpenVPN**: OpenVPN 2.5+ with easy-rsa
- **PPTP**: pptpd daemon
- **L2TP/IPsec**: xl2tpd + strongSwan/Libreswan

## 🔧 Backend Implementation Examples

### Core VPN Manager
```go
// internal/core/manager.go
package core

import (
    "context"
    "sync"
)

type VPNManager struct {
    wireguard *wireguard.Server
    openvpn   *openvpn.Server
    pptp      *pptp.Server
    l2tp      *l2tp.Server

    userManager *UserManager
    statsCollector *StatsCollector

    mu sync.RWMutex
}

func NewVPNManager(cfg *Config) (*VPNManager, error) {
    mgr := &VPNManager{
        userManager: NewUserManager(cfg.DB),
        statsCollector: NewStatsCollector(),
    }

    // Initialize protocol servers
    if cfg.WireGuard.Enabled {
        wg, err := wireguard.NewServer(cfg.WireGuard)
        if err != nil {
            return nil, err
        }
        mgr.wireguard = wg
    }

    // Similar for other protocols...

    return mgr, nil
}

func (m *VPNManager) Start(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Start enabled protocols
    if m.wireguard != nil {
        if err := m.wireguard.Start(); err != nil {
            return err
        }
    }

    // Start stats collection
    go m.statsCollector.Run(ctx)

    return nil
}

func (m *VPNManager) GetStatus() *VPNStatus {
    m.mu.RLock()
    defer m.mu.RUnlock()

    return &VPNStatus{
        WireGuard: m.wireguard.GetStatus(),
        OpenVPN:   m.openvpn.GetStatus(),
        PPTP:      m.pptp.GetStatus(),
        L2TP:      m.l2tp.GetStatus(),
        Stats:     m.statsCollector.GetStats(),
    }
}
```

### WireGuard Implementation
```go
// internal/wireguard/server.go
package wireguard

import (
    "golang.zx2c4.com/wireguard/wgctrl"
    "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Server struct {
    config    *Config
    client    *wgctrl.Client
    device    string
    privateKey wgtypes.Key
    peers     map[string]*Peer
}

func NewServer(cfg *Config) (*Server, error) {
    client, err := wgctrl.New()
    if err != nil {
        return nil, err
    }

    privateKey, err := wgtypes.GeneratePrivateKey()
    if err != nil {
        return nil, err
    }

    return &Server{
        config:     cfg,
        client:     client,
        device:     cfg.Interface,
        privateKey: privateKey,
        peers:      make(map[string]*Peer),
    }, nil
}

func (s *Server) Start() error {
    // Create WireGuard interface
    if err := s.createInterface(); err != nil {
        return err
    }

    // Configure interface
    cfg := wgtypes.Config{
        PrivateKey:   &s.privateKey,
        ListenPort:   &s.config.ListenPort,
        ReplacePeers: false,
    }

    return s.client.ConfigureDevice(s.device, cfg)
}

func (s *Server) AddPeer(peer *Peer) error {
    // Parse peer public key
    pubKey, err := wgtypes.ParseKey(peer.PublicKey)
    if err != nil {
        return err
    }

    // Parse allowed IPs
    allowedIPs, err := parseIPNets(peer.AllowedIPs)
    if err != nil {
        return err
    }

    // Configure peer
    peerCfg := wgtypes.PeerConfig{
        PublicKey:  pubKey,
        AllowedIPs: allowedIPs,
    }

    cfg := wgtypes.Config{
        Peers: []wgtypes.PeerConfig{peerCfg},
    }

    if err := s.client.ConfigureDevice(s.device, cfg); err != nil {
        return err
    }

    s.peers[peer.PublicKey] = peer
    return nil
}

func (s *Server) GetPeerStats() ([]*PeerStats, error) {
    device, err := s.client.Device(s.device)
    if err != nil {
        return nil, err
    }

    var stats []*PeerStats
    for _, peer := range device.Peers {
        stats = append(stats, &PeerStats{
            PublicKey:       peer.PublicKey.String(),
            Endpoint:        peer.Endpoint.String(),
            LastHandshake:   peer.LastHandshakeTime,
            BytesReceived:   peer.ReceiveBytes,
            BytesSent:       peer.TransmitBytes,
        })
    }

    return stats, nil
}

func (s *Server) GenerateClientConfig(peer *Peer) (string, error) {
    config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
DNS = %s

[Peer]
PublicKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25
`, peer.PrivateKey, peer.Address, s.config.DNS,
    s.privateKey.PublicKey().String(), s.config.Endpoint, s.config.ListenPort)

    return config, nil
}
```

### API Handlers
```go
// internal/api/handlers.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    vpnManager *core.VPNManager
}

func NewHandler(mgr *core.VPNManager) *Handler {
    return &Handler{vpnManager: mgr}
}

// GET /api/vpn/status
func (h *Handler) GetStatus(c *gin.Context) {
    status := h.vpnManager.GetStatus()
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    status,
    })
}

// POST /api/vpn/wireguard/peers
func (h *Handler) AddWireGuardPeer(c *gin.Context) {
    var req struct {
        Name       string `json:"name" binding:"required"`
        AllowedIPs string `json:"allowedIPs" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    peer, err := h.vpnManager.WireGuard.CreatePeer(req.Name, req.AllowedIPs)
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

// GET /api/vpn/users
func (h *Handler) GetUsers(c *gin.Context) {
    users, err := h.vpnManager.UserManager.GetAllUsers()
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

// PUT /api/vpn/users/:id/protocols
func (h *Handler) UpdateUserProtocols(c *gin.Context) {
    userID := c.Param("id")

    var req struct {
        WireGuard bool `json:"wireguard"`
        OpenVPN   bool `json:"openvpn"`
        PPTP      bool `json:"pptp"`
        L2TP      bool `json:"l2tp"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    err := h.vpnManager.UserManager.UpdateProtocolAccess(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "User protocols updated successfully",
    })
}

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
        "message": fmt.Sprintf("%s server started successfully", protocol),
    })
}

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
        "message": fmt.Sprintf("%s server stopped successfully", protocol),
    })
}
```

### Database Schema
```sql
-- users table (extends existing user system)
CREATE TABLE vpn_users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_connection TIMESTAMP NULL,
    enabled BOOLEAN DEFAULT true,
    INDEX idx_username (username),
    INDEX idx_enabled (enabled)
);

-- protocol access permissions
CREATE TABLE vpn_user_protocols (
    user_id VARCHAR(36) NOT NULL,
    protocol ENUM('wireguard', 'openvpn', 'pptp', 'l2tp') NOT NULL,
    enabled BOOLEAN DEFAULT false,
    PRIMARY KEY (user_id, protocol),
    FOREIGN KEY (user_id) REFERENCES vpn_users(id) ON DELETE CASCADE
);

-- wireguard peers
CREATE TABLE wireguard_peers (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    name VARCHAR(255) NOT NULL,
    public_key VARCHAR(255) NOT NULL UNIQUE,
    private_key VARCHAR(255) NOT NULL,
    allowed_ips VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES vpn_users(id) ON DELETE SET NULL
);

-- openvpn certificates
CREATE TABLE openvpn_certificates (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    common_name VARCHAR(255) NOT NULL UNIQUE,
    serial_number VARCHAR(255) NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL,
    status ENUM('valid', 'revoked', 'expired') DEFAULT 'valid',
    certificate TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES vpn_users(id) ON DELETE SET NULL
);

-- connection logs
CREATE TABLE vpn_connections (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    protocol ENUM('wireguard', 'openvpn', 'pptp', 'l2tp') NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP NULL,
    bytes_received BIGINT DEFAULT 0,
    bytes_sent BIGINT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES vpn_users(id) ON DELETE CASCADE,
    INDEX idx_user_protocol (user_id, protocol),
    INDEX idx_connected (connected_at)
);

-- server configuration
CREATE TABLE vpn_config (
    protocol ENUM('wireguard', 'openvpn', 'pptp', 'l2tp') PRIMARY KEY,
    config JSON NOT NULL,
    enabled BOOLEAN DEFAULT false,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 🎨 Design Philosophy

The UI follows Stumpfworks design principles:

1. **Gradient Backgrounds**: Purple-to-blue gradients for modern look
2. **Glassmorphism**: Frosted glass effects with backdrop blur
3. **Smooth Animations**: Framer Motion for fluid interactions
4. **Color-Coded Protocols**:
   - WireGuard: Cyan/Blue (modern, fast)
   - OpenVPN: Green (reliable, standard)
   - PPTP: Orange/Red (legacy, warning)
   - L2TP/IPsec: Purple/Pink (enterprise, secure)
5. **Status Indicators**: Live connection pulses, traffic graphs
6. **Responsive Design**: Works on desktop, tablet, and mobile

## 🔒 Security Considerations

1. **PPTP Warning**: Display security warnings for PPTP usage
2. **Strong Defaults**: AES-256 encryption by default
3. **Key Management**: Secure storage of private keys
4. **Certificate Validation**: Proper PKI for OpenVPN
5. **Rate Limiting**: Prevent brute force attacks
6. **Audit Logging**: All configuration changes logged
7. **Firewall Integration**: Automatic iptables rules

## 📊 Performance Targets

- Support 100+ concurrent VPN connections
- <50ms API response times
- Real-time stats updates every 5 seconds
- Minimal CPU overhead (<5% at 50 connections)
- Efficient memory usage (<500MB total)

## 🚢 Deployment

```bash
# Build backend
cd plugins/vpn-server
go build -o vpn-server cmd/vpn-server/main.go

# Install as systemd service
sudo systemctl enable vpn-server
sudo systemctl start vpn-server

# Frontend is bundled with main NAS application
```

## 📝 Configuration Example

```yaml
# /etc/stumpfworks-nas/vpn-server.yaml
wireguard:
  enabled: true
  interface: wg0
  listenPort: 51820
  subnet: 10.8.0.0/24
  dns: 8.8.8.8, 8.8.4.4

openvpn:
  enabled: true
  protocol: udp
  port: 1194
  subnet: 10.9.0.0/24
  cipher: AES-256-GCM

pptp:
  enabled: false
  port: 1723
  subnet: 10.10.0.0/24

l2tp:
  enabled: true
  port: 1701
  ipsecPort: 500
  subnet: 10.11.0.0/24
  psk: "your-pre-shared-key"

general:
  maxConnections: 100
  enableLogging: true
  logLevel: info
  accountSource: local
```

## 🎯 Next Steps

1. **Review this architecture** - Ensure it meets your vision
2. **Backend implementation** - Start with Phase 1 (WireGuard)
3. **Database integration** - Set up schema and migrations
4. **API development** - Implement REST endpoints
5. **Testing** - Unit and integration tests
6. **Documentation** - User guides and API docs
7. **AppStore integration** - Package as installable plugin

This comprehensive VPN Server will make Stumpfworks NAS a compelling choice for users who need robust, beautiful VPN management!
