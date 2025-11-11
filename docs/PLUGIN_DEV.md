# Plugin Development Guide

> Build extensions for Stumpf.Works NAS Solution

---

## Table of Contents

- [Overview](#overview)
- [Plugin Architecture](#plugin-architecture)
- [Plugin Types](#plugin-types)
- [Getting Started](#getting-started)
- [Plugin Manifest](#plugin-manifest)
- [Native Plugins (Go)](#native-plugins-go)
- [Docker Plugins](#docker-plugins)
- [Frontend Plugins](#frontend-plugins)
- [Plugin API](#plugin-api)
- [Security & Sandboxing](#security--sandboxing)
- [Publishing](#publishing)
- [Best Practices](#best-practices)

---

## Overview

The Stumpf.Works plugin system allows developers to extend NAS functionality without modifying the core codebase. Plugins are **first-class citizens** — even core features use the same plugin APIs.

### What Can Plugins Do?

- Add new applications to the Dock
- Provide backend services (APIs, background jobs)
- Integrate third-party services (Plex, Nextcloud, etc.)
- Add system monitoring widgets
- Extend storage or network capabilities
- Automate tasks (scripts, cron jobs)

---

## Plugin Architecture

```
┌─────────────────────────────────────────────────┐
│              Plugin Manager                     │
│  ┌───────────────────────────────────────────┐  │
│  │  Loader  │  Registry  │  Lifecycle        │  │
│  └───────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────┐  │
│  │  Sandbox │  Permissions │  Resource Limits│  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
           │                │               │
   ┌───────▼──────┐  ┌──────▼────┐  ┌──────▼────────┐
   │  Native      │  │  Docker   │  │  Frontend     │
   │  Plugin      │  │  Plugin   │  │  Plugin       │
   │  (Go .so)    │  │  (Container)  │  (React/JS)  │
   └──────────────┘  └───────────┘  └───────────────┘
```

---

## Plugin Types

### 1. Native Plugins (Go)

**Description:** Go modules compiled as shared libraries (`.so` files) and loaded directly into the backend process.

**Pros:**
- Maximum performance
- Direct access to system APIs
- Shared memory with core backend

**Cons:**
- Must be written in Go
- Tightly coupled to backend version
- Potential for crashes (no process isolation)

**Use Cases:**
- High-performance data processing
- Deep system integration
- Extending core functionality

---

### 2. Docker Plugins

**Description:** Containerized services that run in isolated Docker containers.

**Pros:**
- Language-agnostic (Python, Node.js, Rust, etc.)
- Process isolation (cannot crash backend)
- Easy dependency management
- Portable across systems

**Cons:**
- Overhead (container runtime)
- Network latency for API calls
- Requires Docker/Podman

**Use Cases:**
- Third-party services (databases, web apps)
- Services with complex dependencies
- Microservices

---

### 3. Frontend Plugins

**Description:** JavaScript/React components that add UI elements (apps, widgets, panels).

**Pros:**
- Easy to develop (standard React)
- Live in user's browser (no backend needed for simple plugins)
- Can be purely presentational

**Cons:**
- No backend access (unless paired with native/Docker plugin)
- Limited to UI functionality

**Use Cases:**
- Dashboard widgets
- Custom apps
- UI themes

---

### 4. Hybrid Plugins

Combine multiple types:
- Backend (Native or Docker) + Frontend (React)
- Example: A monitoring plugin with a Go backend service and a React dashboard widget

---

## Getting Started

### Prerequisites

- Stumpf.Works NAS installed
- Development environment:
  - Go 1.21+ (for native plugins)
  - Docker (for Docker plugins)
  - Node.js 18+ (for frontend plugins)
- Plugin SDK (provided in `/plugins/sdk/`)

---

### Plugin SDK

The SDK includes:
- Go interfaces and types
- React component stubs
- Example plugins
- Build scripts
- Testing utilities

**Installation:**

```bash
# Clone SDK repository
git clone https://github.com/stumpfworks/plugin-sdk.git
cd plugin-sdk

# Install dependencies
make install
```

---

## Plugin Manifest

Every plugin requires a `plugin.json` file in its root directory.

### Minimal Example

```json
{
  "id": "com.example.hello",
  "name": "Hello World",
  "version": "1.0.0",
  "author": "John Doe",
  "description": "A simple hello world plugin",
  "type": "native",
  "entrypoint": "./hello.so"
}
```

---

### Complete Schema

```json
{
  "id": "com.example.myplugin",
  "name": "My Plugin",
  "version": "1.2.3",
  "author": "John Doe <john@example.com>",
  "description": "A comprehensive plugin example",
  "homepage": "https://github.com/example/myplugin",
  "license": "MIT",

  "type": "native|docker|frontend|hybrid",

  "entrypoint": "./plugin.so",

  "compatibility": {
    "system": ">=1.0.0 <2.0.0"
  },

  "dependencies": {
    "plugins": {
      "com.other.plugin": "^1.0.0"
    },
    "system": {
      "docker": true,
      "zfs": false
    }
  },

  "permissions": [
    "storage.read",
    "storage.write",
    "network.read",
    "system.restart"
  ],

  "api": {
    "prefix": "/plugins/myplugin",
    "endpoints": [
      {
        "method": "GET",
        "path": "/status",
        "handler": "GetStatus",
        "auth": true
      },
      {
        "method": "POST",
        "path": "/action",
        "handler": "PerformAction",
        "auth": true
      }
    ]
  },

  "frontend": {
    "entry": "./dist/index.js",
    "icon": "./icon.svg",
    "window": {
      "title": "My Plugin",
      "width": 800,
      "height": 600,
      "resizable": true,
      "minWidth": 400,
      "minHeight": 300
    }
  },

  "docker": {
    "image": "example/myplugin:1.2.3",
    "ports": {
      "8080": "8080"
    },
    "volumes": {
      "/data": "/var/lib/myplugin"
    },
    "environment": {
      "API_KEY": "${PLUGIN_API_KEY}"
    }
  },

  "settings": {
    "schema": {
      "type": "object",
      "properties": {
        "apiKey": {
          "type": "string",
          "title": "API Key",
          "description": "Your API key"
        },
        "enabled": {
          "type": "boolean",
          "title": "Enabled",
          "default": true
        }
      }
    }
  },

  "install": {
    "script": "./install.sh"
  },

  "uninstall": {
    "script": "./uninstall.sh"
  }
}
```

---

## Native Plugins (Go)

### Plugin Interface

All native plugins must implement the `Plugin` interface:

```go
// plugin.go
package main

import (
    "github.com/stumpfworks/nas/sdk/plugin"
)

type MyPlugin struct {
    config *plugin.Config
}

func (p *MyPlugin) Init(config *plugin.Config) error {
    p.config = config
    // Initialize plugin (load config, connect to DB, etc.)
    return nil
}

func (p *MyPlugin) Start() error {
    // Start background services, goroutines, etc.
    return nil
}

func (p *MyPlugin) Stop() error {
    // Clean shutdown (close connections, save state, etc.)
    return nil
}

func (p *MyPlugin) GetInfo() plugin.Info {
    return plugin.Info{
        ID:          "com.example.myplugin",
        Name:        "My Plugin",
        Version:     "1.0.0",
        Author:      "John Doe",
        Description: "A sample plugin",
    }
}

func (p *MyPlugin) RegisterRoutes(router plugin.Router) {
    router.GET("/status", p.GetStatus)
    router.POST("/action", p.PerformAction)
}

func (p *MyPlugin) GetStatus(ctx plugin.Context) error {
    return ctx.JSON(200, map[string]string{
        "status": "ok",
    })
}

func (p *MyPlugin) PerformAction(ctx plugin.Context) error {
    var req struct {
        Action string `json:"action"`
    }
    if err := ctx.Bind(&req); err != nil {
        return err
    }

    // Perform action...

    return ctx.JSON(200, map[string]string{
        "result": "success",
    })
}

// Required: exported plugin constructor
func NewPlugin() plugin.Plugin {
    return &MyPlugin{}
}
```

---

### Building a Native Plugin

```bash
# Build as shared library
go build -buildmode=plugin -o myplugin.so plugin.go

# Create plugin package
mkdir -p dist
cp myplugin.so dist/
cp plugin.json dist/
cp icon.svg dist/

# Package as .tar.gz
tar -czf myplugin-1.0.0.tar.gz -C dist .
```

---

### Plugin SDK API

#### Storage Access

```go
import "github.com/stumpfworks/nas/sdk/storage"

// Read file
data, err := p.config.Storage.ReadFile("/data/config.json")

// Write file
err := p.config.Storage.WriteFile("/data/config.json", data, 0644)

// List volumes
volumes, err := p.config.Storage.ListVolumes()
```

#### Database Access

```go
import "github.com/stumpfworks/nas/sdk/database"

// Query
var users []User
err := p.config.DB.Find(&users).Error

// Insert
user := User{Name: "John"}
err := p.config.DB.Create(&user).Error
```

#### System Commands

```go
import "github.com/stumpfworks/nas/sdk/system"

// Execute command
output, err := p.config.System.Exec("df", "-h")

// Restart service
err := p.config.System.RestartService("smbd")
```

#### WebSocket Events

```go
import "github.com/stumpfworks/nas/sdk/events"

// Publish event
p.config.Events.Publish("plugin.myplugin.status", map[string]interface{}{
    "status": "running",
})

// Subscribe to events
p.config.Events.Subscribe("system.storage.changed", func(data interface{}) {
    // Handle event
})
```

---

## Docker Plugins

### Dockerfile Example

```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 8080

CMD ["node", "server.js"]
```

### Server Example (Node.js)

```javascript
// server.js
const express = require('express');
const app = express();

app.use(express.json());

app.get('/status', (req, res) => {
  res.json({ status: 'ok' });
});

app.post('/action', (req, res) => {
  const { action } = req.body;
  // Perform action...
  res.json({ result: 'success' });
});

app.listen(8080, () => {
  console.log('Plugin running on port 8080');
});
```

### plugin.json for Docker Plugin

```json
{
  "id": "com.example.dockerplugin",
  "name": "Docker Plugin",
  "version": "1.0.0",
  "type": "docker",
  "docker": {
    "image": "example/myplugin:1.0.0",
    "ports": {
      "8080": "8080"
    },
    "environment": {
      "API_URL": "${SYSTEM_API_URL}",
      "API_TOKEN": "${PLUGIN_API_TOKEN}"
    }
  },
  "api": {
    "proxy": {
      "prefix": "/plugins/dockerplugin",
      "target": "http://localhost:8080"
    }
  }
}
```

---

## Frontend Plugins

### React Component Example

```tsx
// MyPluginApp.tsx
import React, { useState, useEffect } from 'react';
import { usePlugin } from '@stumpfworks/plugin-sdk';

export default function MyPluginApp() {
  const { api, config } = usePlugin();
  const [status, setStatus] = useState('loading');

  useEffect(() => {
    api.get('/plugins/myplugin/status')
      .then(res => setStatus(res.data.status))
      .catch(err => setStatus('error'));
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">My Plugin</h1>
      <p>Status: {status}</p>
      <button
        className="btn-primary"
        onClick={() => {
          api.post('/plugins/myplugin/action', { action: 'test' });
        }}
      >
        Perform Action
      </button>
    </div>
  );
}
```

### Building Frontend Plugin

```bash
# Install dependencies
npm install

# Build (Vite or Webpack)
npm run build

# Output: dist/index.js
```

---

## Plugin API

### REST API

All plugins can register custom API endpoints under `/api/v1/plugins/{plugin-id}`.

**Example:**
- Plugin ID: `com.example.myplugin`
- Endpoint: `/status`
- Full URL: `/api/v1/plugins/com.example.myplugin/status`

### WebSocket Events

Plugins can publish and subscribe to events via WebSocket.

**Event Channels:**
- `plugin.{id}.{event}` - Plugin-specific events
- `system.*` - System events (requires permission)

---

## Security & Sandboxing

### Permission System

Plugins declare required permissions in `plugin.json`:

```json
{
  "permissions": [
    "storage.read",       // Read storage information
    "storage.write",      // Modify storage (create volumes, etc.)
    "network.read",       // View network configuration
    "network.write",      // Modify network settings
    "users.read",         // View user list
    "users.write",        // Create/modify users
    "system.read",        // View system info (CPU, RAM, etc.)
    "system.write",       // Restart services, shutdown
    "docker.read",        // List containers
    "docker.write",       // Create/manage containers
    "files.read",         // Read files
    "files.write"         // Write files
  ]
}
```

**Permission Enforcement:**
- Checked at runtime by Plugin Manager
- Unauthorized actions return `403 Forbidden`
- User can review permissions before installation

---

### Sandboxing

#### Native Plugins
- Run in same process (shared memory)
- Limited sandboxing (trust-based)
- Code review required for official plugins

#### Docker Plugins
- Full process isolation
- Network isolation (custom bridge)
- Resource limits (CPU, memory)
- Read-only filesystem (except mounted volumes)

---

### Resource Limits

Plugins can be limited:

```json
{
  "resources": {
    "cpu": "0.5",      // Max 50% of one CPU core
    "memory": "256M",  // Max 256MB RAM
    "disk": "1G"       // Max 1GB storage
  }
}
```

---

## Publishing

### Plugin Marketplace

Official plugins are hosted in the Stumpf.Works Plugin Marketplace.

**Submission Process:**

1. **Develop Plugin** - Follow this guide
2. **Test Locally** - Install via "Upload Plugin" in Plugin Center
3. **Create GitHub Repo** - Public repository with source code
4. **Submit PR** - Add plugin to [marketplace registry](https://github.com/stumpfworks/plugin-registry)
5. **Code Review** - Stumpf.Works team reviews code
6. **Approval** - Plugin appears in marketplace

---

### Plugin Package Format

```
myplugin-1.0.0.tar.gz
├── plugin.json       # Manifest
├── icon.svg          # Plugin icon (512x512)
├── README.md         # Documentation
├── LICENSE           # License file
├── plugin.so         # Native plugin binary (if native)
├── dist/             # Frontend assets (if frontend)
│   └── index.js
└── docker-compose.yml  # Docker config (if Docker)
```

---

### Versioning

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward-compatible)
- **PATCH**: Bug fixes

Example: `1.2.3`

---

## Best Practices

### 1. Error Handling

Always handle errors gracefully:

```go
func (p *MyPlugin) GetData(ctx plugin.Context) error {
    data, err := fetchData()
    if err != nil {
        return ctx.JSON(500, map[string]string{
            "error": "Failed to fetch data",
        })
    }
    return ctx.JSON(200, data)
}
```

---

### 2. Logging

Use structured logging:

```go
p.config.Logger.Info("Plugin started", "version", "1.0.0")
p.config.Logger.Error("Failed to connect", "error", err)
```

---

### 3. Configuration

Store plugin config in `/etc/stumpfworks/plugins/{id}/config.json`:

```go
func (p *MyPlugin) LoadConfig() error {
    data, err := p.config.Storage.ReadFile("/etc/stumpfworks/plugins/myplugin/config.json")
    if err != nil {
        return err
    }
    return json.Unmarshal(data, &p.settings)
}
```

---

### 4. Testing

Write unit tests:

```go
func TestGetStatus(t *testing.T) {
    plugin := &MyPlugin{}
    plugin.Init(mockConfig())

    resp := plugin.GetStatus(mockContext())

    assert.Equal(t, 200, resp.StatusCode)
}
```

---

### 5. Documentation

Include comprehensive README.md:
- Installation instructions
- Configuration options
- API endpoints
- Screenshots
- Troubleshooting

---

### 6. Performance

- Avoid blocking operations in API handlers
- Use goroutines for background tasks
- Cache expensive operations
- Monitor resource usage

---

### 7. Security

- Validate all user input
- Sanitize file paths (prevent directory traversal)
- Use prepared statements for SQL queries
- Never log sensitive data (passwords, tokens)
- Keep dependencies updated

---

## Example Plugins

### 1. System Monitor Widget

**Type:** Frontend
**Description:** Dashboard widget showing CPU/RAM graphs

### 2. Backup to S3

**Type:** Native (Go)
**Description:** Automated backups to AWS S3

### 3. Plex Integration

**Type:** Docker
**Description:** Embeds Plex Media Server

### 4. Custom Dashboard

**Type:** Hybrid (Go + React)
**Description:** Custom monitoring dashboard with backend API

---

## Support & Resources

- **Plugin SDK:** [github.com/stumpfworks/plugin-sdk](https://github.com/stumpfworks/plugin-sdk)
- **Examples:** [github.com/stumpfworks/plugin-examples](https://github.com/stumpfworks/plugin-examples)
- **Forum:** [community.stumpf.works](https://community.stumpf.works)
- **Discord:** [discord.gg/stumpfworks](https://discord.gg/stumpfworks)

---

**Document Version:** 1.0.0
**Last Updated:** 2025-11-11
**Plugin API Version:** v1
