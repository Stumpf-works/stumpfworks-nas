# System Architecture

> Comprehensive architectural design for Stumpf.Works NAS Solution

---

## Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [System Layers](#system-layers)
- [Component Hierarchy](#component-hierarchy)
- [Backend Architecture](#backend-architecture)
- [Frontend Architecture](#frontend-architecture)
- [Plugin System](#plugin-system)
- [Data Flow](#data-flow)
- [Security Architecture](#security-architecture)
- [Deployment Architecture](#deployment-architecture)

---

## Overview

Stumpf.Works NAS is a **layered, modular system** built on Debian, designed for extensibility, security, and performance. The architecture follows modern microservices principles while maintaining simplicity for a single-node NAS deployment.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        End Users                            │
│              (Browser, Mobile App, API Clients)             │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Frontend Layer                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  React SPA (Vite + TailwindCSS + Framer Motion)     │   │
│  │  • Desktop Environment  • Window Manager             │   │
│  │  • Dock • Apps • Components                          │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS + WebSocket
┌────────────────────────▼────────────────────────────────────┐
│                     API Gateway                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  REST API + WebSocket Server (Go)                    │   │
│  │  • Authentication • Rate Limiting • CORS             │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Backend Core (Go)                         │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐   │
│  │ Storage  │ Network  │ Users    │ System   │ Plugins  │   │
│  │ Manager  │ Manager  │ Manager  │ Monitor  │ Manager  │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘   │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐   │
│  │ Share    │ Container│ VM       │ Backup   │ Update   │   │
│  │ Manager  │ Manager  │ Manager  │ Manager  │ Service  │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  System Layer (Debian)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  LVM • mdadm • ZFS • Docker • libvirt • systemd      │   │
│  │  SMB • NFS • FTP • Firewall • Network Stack          │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  Hardware Layer                             │
│           CPU • RAM • Disks • Network Interfaces            │
└─────────────────────────────────────────────────────────────┘
```

---

## Design Principles

### 1. Modularity
Every component is independently developed, tested, and deployed. Modules communicate through well-defined interfaces.

### 2. Separation of Concerns
- **Frontend:** UI/UX only, no business logic
- **Backend:** Business logic, data processing, system integration
- **System Layer:** Low-level OS operations

### 3. Plugin-First Architecture
Core functionality is built using the same plugin APIs available to third-party developers.

### 4. Security by Default
- Principle of least privilege
- Encrypted communication (HTTPS/TLS)
- Sandboxed plugin execution
- Immutable audit logs

### 5. Performance
- Async operations for I/O-bound tasks
- Caching strategies (in-memory, Redis optional)
- Lazy loading for frontend components
- Database query optimization

### 6. Testability
- Dependency injection
- Mock-friendly interfaces
- Comprehensive test coverage
- Isolated test environments

---

## System Layers

### Layer 1: Hardware
Physical or virtual hardware resources.

### Layer 2: Operating System (Debian Bookworm)
- Kernel 6.1+
- systemd init system
- Standard Debian packages
- Custom repository for Stumpf.Works packages

### Layer 3: System Services
Low-level services managed by the backend:
- Storage subsystems (LVM, mdadm, ZFS)
- Network stack (interfaces, routing, firewall)
- File sharing (Samba, NFS, FTP)
- Containerization (Docker/Podman)
- Virtualization (libvirt/KVM)

### Layer 4: Backend Core
Go-based application server:
- REST API
- WebSocket server
- Plugin runtime
- System management logic
- Database layer

### Layer 5: Frontend
React-based web application:
- macOS-like UI
- Window management
- Application ecosystem
- Real-time updates

### Layer 6: Plugins
Extend functionality through:
- Native Go modules
- Docker containers
- Frontend components

---

## Component Hierarchy

### Backend Component Tree

```
stumpfworks-backend/
│
├── cmd/
│   └── stumpfworks-server/
│       └── main.go                 # Entry point
│
├── internal/                       # Private application code
│   ├── api/
│   │   ├── router.go               # API routing
│   │   ├── middleware/             # Auth, CORS, logging
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── storage.go
│   │   │   ├── network.go
│   │   │   ├── users.go
│   │   │   ├── system.go
│   │   │   └── plugins.go
│   │   └── websocket/              # WebSocket server
│   │
│   ├── storage/                    # Storage management
│   │   ├── manager.go
│   │   ├── lvm.go
│   │   ├── mdadm.go
│   │   ├── zfs.go                  # Optional
│   │   ├── smart.go
│   │   └── filesystem.go
│   │
│   ├── network/                    # Network management
│   │   ├── manager.go
│   │   ├── interfaces.go
│   │   ├── firewall.go             # ufw integration
│   │   ├── routing.go
│   │   └── vlan.go
│   │
│   ├── users/                      # User management
│   │   ├── manager.go
│   │   ├── authentication.go       # JWT, sessions
│   │   ├── authorization.go        # RBAC
│   │   ├── ldap.go                 # LDAP/AD integration
│   │   └── quota.go
│   │
│   ├── system/                     # System monitoring
│   │   ├── info.go                 # CPU, RAM, uptime
│   │   ├── metrics.go              # Prometheus metrics
│   │   ├── logs.go                 # Log aggregation
│   │   └── updates.go
│   │
│   ├── plugins/                    # Plugin system
│   │   ├── manager.go              # Plugin lifecycle
│   │   ├── loader.go               # Dynamic loading
│   │   ├── api.go                  # Plugin API interface
│   │   ├── sandbox.go              # Security isolation
│   │   └── registry.go             # Plugin repository
│   │
│   ├── shares/                     # File sharing
│   │   ├── smb.go                  # Samba configuration
│   │   ├── nfs.go
│   │   ├── ftp.go
│   │   └── webdav.go
│   │
│   ├── containers/                 # Container management
│   │   ├── docker.go
│   │   ├── podman.go
│   │   └── compose.go
│   │
│   ├── vms/                        # Virtual machines
│   │   ├── libvirt.go
│   │   ├── qemu.go
│   │   └── vnc.go
│   │
│   ├── backup/                     # Backup systems
│   │   ├── manager.go
│   │   ├── snapshot.go
│   │   ├── replication.go
│   │   └── cloud_sync.go
│   │
│   ├── database/                   # Database layer
│   │   ├── db.go                   # SQLite/PostgreSQL
│   │   ├── migrations/
│   │   └── models/
│   │
│   └── config/                     # Configuration
│       ├── config.go
│       └── validation.go
│
├── pkg/                            # Public libraries
│   ├── logger/                     # Structured logging
│   ├── errors/                     # Error handling
│   ├── validator/                  # Input validation
│   └── utils/                      # Helper functions
│
└── go.mod
```

---

## Frontend Architecture

### Component Tree

```
frontend/
│
├── src/
│   ├── main.tsx                    # Entry point
│   ├── App.tsx                     # Root component
│   │
│   ├── layout/                     # Core UI layout
│   │   ├── Desktop.tsx             # Desktop environment
│   │   ├── Dock.tsx                # Bottom dock
│   │   ├── TopBar.tsx              # Menu bar
│   │   ├── ControlCenter.tsx       # Quick settings
│   │   ├── Launchpad.tsx           # App grid
│   │   ├── NotificationCenter.tsx
│   │   └── WindowManager.tsx       # Multi-window system
│   │
│   ├── apps/                       # Applications
│   │   ├── Dashboard/
│   │   │   ├── Dashboard.tsx
│   │   │   ├── components/
│   │   │   └── hooks/
│   │   ├── StorageManager/
│   │   ├── FileStation/
│   │   ├── UserManager/
│   │   ├── NetworkManager/
│   │   ├── ShareManager/
│   │   ├── ContainerManager/
│   │   ├── PluginCenter/
│   │   └── Settings/
│   │
│   ├── components/                 # Shared components
│   │   ├── ui/                     # UI primitives
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── ...
│   │   ├── Window/                 # Window component
│   │   └── animations/             # Framer Motion wrappers
│   │
│   ├── store/                      # State management
│   │   ├── index.ts                # Zustand store
│   │   ├── slices/
│   │   │   ├── authSlice.ts
│   │   │   ├── windowSlice.ts
│   │   │   ├── systemSlice.ts
│   │   │   └── pluginSlice.ts
│   │   └── middleware/
│   │
│   ├── api/                        # API client
│   │   ├── client.ts               # Axios instance
│   │   ├── auth.ts
│   │   ├── storage.ts
│   │   ├── network.ts
│   │   ├── users.ts
│   │   └── websocket.ts            # WebSocket client
│   │
│   ├── hooks/                      # Custom React hooks
│   │   ├── useAuth.ts
│   │   ├── useWebSocket.ts
│   │   ├── useSystemMetrics.ts
│   │   └── usePlugins.ts
│   │
│   ├── utils/                      # Helper functions
│   ├── types/                      # TypeScript types
│   ├── styles/                     # Global styles
│   └── assets/                     # Images, icons
│
├── public/
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

---

## Backend Architecture

### API Design

#### REST API

**Base URL:** `/api/v1`

**Endpoints:**

```
Authentication & Sessions
POST   /auth/login
POST   /auth/logout
POST   /auth/refresh
GET    /auth/me

System Information
GET    /system/info              # CPU, RAM, uptime
GET    /system/metrics           # Real-time metrics
GET    /system/logs              # System logs
POST   /system/restart
POST   /system/shutdown

Storage Management
GET    /storage/disks            # List all disks
GET    /storage/disks/:id        # Disk details
GET    /storage/volumes          # List volumes
POST   /storage/volumes          # Create volume (LVM/mdadm)
DELETE /storage/volumes/:id
GET    /storage/smart/:disk      # SMART data
GET    /storage/snapshots
POST   /storage/snapshots        # Create snapshot

Network Management
GET    /network/interfaces
POST   /network/interfaces       # Configure interface
GET    /network/firewall
POST   /network/firewall         # Add firewall rule
GET    /network/routes

User Management
GET    /users
POST   /users                    # Create user
PUT    /users/:id
DELETE /users/:id
GET    /users/:id/permissions

File Shares
GET    /shares                   # List all shares
POST   /shares                   # Create share (SMB/NFS)
PUT    /shares/:id
DELETE /shares/:id

Plugins
GET    /plugins                  # List installed plugins
POST   /plugins/install          # Install plugin
POST   /plugins/:id/enable
POST   /plugins/:id/disable
DELETE /plugins/:id              # Uninstall
GET    /plugins/marketplace      # Available plugins

Containers (Docker/Podman)
GET    /containers
POST   /containers/create
POST   /containers/:id/start
POST   /containers/:id/stop
DELETE /containers/:id

Virtual Machines
GET    /vms
POST   /vms/create
POST   /vms/:id/start
POST   /vms/:id/stop
GET    /vms/:id/vnc              # VNC console URL
```

#### WebSocket API

**Endpoint:** `/ws`

**Message Types:**

```json
{
  "type": "subscribe",
  "channel": "system.metrics"
}

{
  "type": "event",
  "channel": "system.metrics",
  "data": {
    "cpu": 45.2,
    "memory": 67.8,
    "timestamp": 1699999999
  }
}
```

**Channels:**
- `system.metrics` - Real-time system metrics
- `storage.events` - Storage changes
- `network.events` - Network state changes
- `plugins.events` - Plugin lifecycle events
- `logs.stream` - Live log streaming

---

## Plugin System

### Plugin Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Plugin Manager                     │
├─────────────────────────────────────────────────────┤
│  Loader │ Registry │ Lifecycle │ Sandbox │ API     │
└─────────────────────────────────────────────────────┘
           │              │                │
   ┌───────▼────┐  ┌──────▼──────┐  ┌──────▼──────┐
   │   Native   │  │   Docker    │  │  Frontend   │
   │   Plugin   │  │   Plugin    │  │   Plugin    │
   └────────────┘  └─────────────┘  └─────────────┘
```

### Plugin Manifest (`plugin.json`)

```json
{
  "id": "com.example.myplugin",
  "name": "My Plugin",
  "version": "1.0.0",
  "author": "John Doe",
  "description": "A sample plugin",
  "type": "native|docker|frontend",
  "entrypoint": "./plugin.so",
  "dependencies": {
    "system": ">=1.0.0",
    "plugins": ["com.other.plugin"]
  },
  "permissions": [
    "storage.read",
    "network.write"
  ],
  "api": {
    "endpoints": [
      {
        "method": "GET",
        "path": "/myplugin/data",
        "handler": "GetData"
      }
    ]
  },
  "frontend": {
    "entry": "./dist/index.js",
    "icon": "./icon.svg",
    "window": {
      "title": "My Plugin",
      "width": 800,
      "height": 600
    }
  }
}
```

### Plugin API Interface

```go
type Plugin interface {
    Init() error
    Start() error
    Stop() error
    GetInfo() PluginInfo
    RegisterRoutes(router *Router)
}

type PluginInfo struct {
    ID          string
    Name        string
    Version     string
    Author      string
    Description string
}
```

---

## Data Flow

### 1. User Authentication Flow

```
User → Frontend → POST /auth/login → Backend
                                    ↓
                              Validate credentials
                                    ↓
                              Generate JWT token
                                    ↓
Frontend ← JSON response ← Backend (token + user info)
    ↓
Store token in localStorage
    ↓
Include in Authorization header for all requests
```

### 2. Real-Time Metrics Flow

```
Backend System Monitor (goroutine)
    ↓
Collect metrics every 1s
    ↓
Publish to WebSocket channel
    ↓
Frontend (subscribed clients)
    ↓
Update UI components (charts, graphs)
```

### 3. Storage Volume Creation Flow

```
Frontend → POST /storage/volumes → Backend
                                      ↓
                              Validate request
                                      ↓
                              Check permissions
                                      ↓
                          Execute system commands
                          (lvcreate, mdadm, etc.)
                                      ↓
                          Update database
                                      ↓
                          Emit event via WebSocket
                                      ↓
Frontend ← Response ← Backend (success/error)
    ↓
Update UI (show new volume)
```

---

## Security Architecture

### 1. Authentication
- **JWT tokens** (access + refresh)
- **Session management** (Redis or in-memory)
- **Password hashing** (bcrypt)
- **2FA support** (TOTP)

### 2. Authorization
- **Role-Based Access Control (RBAC)**
  - Roles: Admin, User, Guest
  - Permissions: read, write, execute per resource
- **API key support** for programmatic access

### 3. Network Security
- **HTTPS only** (Let's Encrypt or self-signed)
- **CORS policies**
- **Rate limiting** (per IP, per user)
- **Firewall integration** (ufw)

### 4. Plugin Security
- **Sandboxed execution** (separate processes or containers)
- **Permission system** (plugins request specific permissions)
- **Code signing** (verify plugin integrity)
- **Resource limits** (CPU, memory quotas)

### 5. Data Security
- **Encrypted storage** (LUKS for volumes)
- **Encrypted communication** (TLS 1.3)
- **Audit logging** (immutable logs)
- **Secure defaults** (minimal permissions)

---

## Deployment Architecture

### Single-Node Deployment (Primary Use Case)

```
┌─────────────────────────────────────────────────┐
│              Debian Host (NAS)                  │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  nginx (reverse proxy)                    │  │
│  │  ├─ HTTPS termination                     │  │
│  │  ├─ /api → backend:8080                   │  │
│  │  └─ / → frontend static files             │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  stumpfworks-backend (Go binary)          │  │
│  │  systemd service                          │  │
│  │  Port: 8080 (internal)                    │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  SQLite database                          │  │
│  │  /var/lib/stumpfworks/db.sqlite           │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  Plugins                                  │  │
│  │  /var/lib/stumpfworks/plugins/            │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

### Docker Deployment (Development/Testing)

```yaml
# docker-compose.yml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./data:/data
    environment:
      - DB_PATH=/data/stumpfworks.db

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend
```

---

## Performance Considerations

### Backend
- **Goroutines** for concurrent operations
- **Connection pooling** (database, HTTP clients)
- **Caching layer** (in-memory or Redis)
- **Lazy loading** for expensive operations
- **Background jobs** for long-running tasks

### Frontend
- **Code splitting** (route-based lazy loading)
- **Virtual scrolling** (large lists)
- **Debounced API calls**
- **Optimistic UI updates**
- **Service workers** (offline support)

### Database
- **Indexed queries**
- **Migration versioning**
- **Connection limits**
- **Query optimization**

---

## Monitoring & Observability

### Metrics (Prometheus)
- System metrics (CPU, RAM, disk I/O)
- API request rate and latency
- Plugin execution time
- Error rates

### Logging
- Structured JSON logs
- Log levels (debug, info, warn, error)
- Log rotation
- Centralized log aggregation (optional)

### Tracing (Optional)
- OpenTelemetry integration
- Request tracing across services

---

## Disaster Recovery

### Backup Strategy
- **Configuration backup** (database, config files)
- **Automated snapshots**
- **Off-site replication**
- **Recovery documentation**

### High Availability (Future)
- Clustered deployment
- Shared storage (NFS, Ceph)
- Load balancing

---

## Technology Stack Summary

| Layer | Technology |
|-------|-----------|
| **OS** | Debian Bookworm (Stable) |
| **Backend** | Go 1.21+ |
| **Frontend** | React 18 + TypeScript |
| **Styling** | TailwindCSS 3 |
| **Animations** | Framer Motion |
| **Build Tool** | Vite |
| **Database** | SQLite (default), PostgreSQL (optional) |
| **API** | REST + WebSocket |
| **Web Server** | nginx (reverse proxy) |
| **Containers** | Docker / Podman |
| **VMs** | libvirt + KVM |
| **Storage** | LVM, mdadm, ZFS (optional) |

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-11
**Status:** Living Document
