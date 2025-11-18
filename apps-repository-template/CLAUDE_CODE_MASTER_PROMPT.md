# StumpfWorks NAS Plugin Development - Master Prompt for Claude Code

**Version**: 1.0.0
**Last Updated**: 2024-12-01
**Purpose**: Comprehensive context for Claude Code when developing StumpfWorks NAS plugins

---

## 🎯 Session Goals

When you start a Claude Code session for StumpfWorks NAS plugin development, you should:

1. **Understand** the StumpfWorks NAS architecture
2. **Follow** established patterns and best practices
3. **Create** production-ready, secure plugins
4. **Document** thoroughly
5. **Test** comprehensively

---

## 📚 StumpfWorks NAS Architecture Overview

### Core System

**StumpfWorks NAS** is a modern, open-source Network-Attached Storage operating system built with:

- **Backend**: Go 1.21+
- **Frontend**: React 18 + TypeScript + TailwindCSS
- **Database**: SQLite with GORM ORM
- **UI Style**: macOS-inspired windowed interface
- **Authentication**: JWT + 2FA support

### System Library Architecture

The core system provides a **System Library** (`backend/internal/system/lib.go`) that offers unified access to:

```go
type SystemLibrary struct {
    Shell   *ShellExecutor    // Safe command execution
    Metrics *MetricsCollector // System metrics
    Storage *StorageManager   // ZFS, RAID, BTRFS, LVM, SMART
    Network *NetworkManager   // Interfaces, firewall, DNS
    Sharing *SharingManager   // Samba, NFS, iSCSI, WebDAV, FTP
    Users   *UserManager      // User management
}
```

**Key Principles**:
- **Singleton Pattern**: Global instance via `system.Get()`
- **Graceful Degradation**: Services fail independently
- **Thread-Safe**: All managers use mutex locks
- **Shell Executor**: All system commands go through centralized executor

### API Structure

**REST API** built with Chi router:
- Base URL: `/api/v1`
- Authentication: JWT Bearer tokens
- Middleware: Logger, Auth, CORS, Admin-only
- Error Handling: Structured error responses

**Common Patterns**:
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    // 1. Get System Library
    lib := system.Get()

    // 2. Check subsystem availability
    if lib.Storage == nil {
        utils.RespondError(w, errors.BadRequest("Service not available", nil))
        return
    }

    // 3. Call service method
    result, err := lib.Storage.Method()

    // 4. Handle errors
    if err != nil {
        utils.RespondError(w, errors.InternalServerError("Error", err))
        return
    }

    // 5. Return success
    utils.RespondSuccess(w, result)
}
```

---

## 🔌 Plugin System Architecture

### Plugin Types

1. **Standalone Go Plugin**
   - Single binary
   - No external dependencies
   - Direct system integration

2. **Docker-Based Plugin**
   - Runs in containers
   - Isolated environment
   - Uses `docker-compose.yml`

3. **Full-Stack Plugin**
   - Go backend + React frontend
   - REST API + Web UI
   - Complete feature set

### Plugin Lifecycle

```
Install → Configure → Enable → Start → Running → Stop → Disable → Uninstall
```

### Plugin Manifest (plugin.json)

**Required Structure**:
```json
{
  "id": "com.company.plugin-name",          // Unique ID (reverse domain)
  "name": "Plugin Display Name",             // Human-readable name
  "version": "1.0.0",                        // Semantic versioning
  "author": "Author Name",                   // Plugin author
  "description": "Short description",        // 50-150 characters
  "icon": "🔌",                              // Emoji or icon path
  "category": "utilities",                   // See categories below
  "entryPoint": "plugin-binary",             // Executable name

  "requires": {
    "docker": false,                         // Needs Docker?
    "ports": [8080, 8443],                   // Required ports
    "storage": "1GB",                        // Minimum storage
    "minNasVersion": "0.1.0"                 // Minimum NAS version
  },

  "config": {                                // Default configuration
    "setting1": "value1"
  },

  "permissions": [                           // Required permissions
    "storage.read",
    "network.write"
  ]
}
```

### Environment Variables

Plugins receive these environment variables:

```bash
PLUGIN_ID=com.company.plugin-name
PLUGIN_DIR=/var/lib/stumpfworks/plugins/plugin-name
NAS_API_URL=http://localhost:8080/api/v1
NAS_API_TOKEN=<auth-token>
```

### Plugin Registry System

**External Repository**: https://github.com/Stumpf-works/stumpfworks-nas-apps

**Registry URL**:
```
https://raw.githubusercontent.com/Stumpf-works/stumpfworks-nas-apps/main/registry.json
```

**Installation Flow**:
1. User browses Plugin Store in UI
2. StumpfWorks NAS queries registry.json
3. User clicks "Install"
4. Plugin downloaded from GitHub Releases
5. Extracted to `/var/lib/stumpfworks/plugins/`
6. Record saved in database
7. Plugin can be enabled/started

---

## 🏗️ Plugin Development Workflow

### 1. Planning Phase

**Before coding, answer these questions**:
- What problem does this solve?
- Does a similar plugin exist?
- What are the system requirements?
- Does it need Docker or can it run natively?
- What permissions/ports are needed?
- How will users configure it?

**Research Existing Patterns**:
- Look at similar plugins in `/plugins/`
- Check StumpfWorks NAS documentation
- Review System Library capabilities

### 2. Setup Phase

**Choose Template**:
```bash
# In stumpfworks-nas-apps repository
cd plugins/
cp -r ../templates/docker-plugin my-plugin/
cd my-plugin/
```

**Initialize Git** (if separate repo):
```bash
git init
git remote add origin https://github.com/username/my-plugin.git
```

**Edit plugin.json**:
- Set unique ID (com.github-username.plugin-name)
- Define requirements
- Set proper category

### 3. Development Phase

**Backend Development**:

```go
// main.go structure
package main

import (
    "os"
    "github.com/rs/zerolog/log"
)

func main() {
    // 1. Setup logging
    setupLogging()

    // 2. Load configuration
    cfg := loadConfig()

    // 3. Initialize services
    service := initializeService(cfg)

    // 4. Start HTTP server
    server := startServer(service)

    // 5. Graceful shutdown
    waitForShutdown(server)
}
```

**API Server Pattern**:
```go
// Use Chi router
r := chi.NewRouter()

// Middleware
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(cors.Handler(corsOptions))

// Health check (required!)
r.Get("/health", healthCheck)

// API routes
r.Route("/api/v1", func(r chi.Router) {
    r.Get("/resource", listResources)
    r.Post("/resource", createResource)
})

// Start server
http.ListenAndServe(":8080", r)
```

**Docker Integration**:
```yaml
# docker-compose.yml
version: '3.8'

services:
  plugin-service:
    image: plugin:latest
    build: ./
    ports:
      - "8080:8080"
    environment:
      - PLUGIN_ID=${PLUGIN_ID}
      - NAS_API_URL=${NAS_API_URL}
    volumes:
      - plugin-data:/data
    restart: unless-stopped

volumes:
  plugin-data:
```

**Frontend Development** (if needed):

```typescript
// React component structure
import React from 'react';

export const PluginDashboard: React.FC = () => {
  const [data, setData] = useState([]);

  useEffect(() => {
    // Fetch data from plugin API
    fetch('http://localhost:8080/api/v1/resource')
      .then(res => res.json())
      .then(data => setData(data));
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Plugin Dashboard</h1>
      {/* Content */}
    </div>
  );
};
```

### 4. Testing Phase

**Local Testing**:
```bash
# Build plugin
./build.sh

# Test with Docker
docker-compose up -d

# Check logs
docker-compose logs -f

# Test API
curl http://localhost:8080/health
```

**Integration Testing**:
```bash
# Copy to StumpfWorks NAS
sudo cp -r . /var/lib/stumpfworks/plugins/my-plugin/

# Enable plugin
curl -X POST http://nas-ip:8080/api/v1/plugins/my-plugin/enable \
  -H "Authorization: Bearer $TOKEN"

# Check status
curl http://nas-ip:8080/api/v1/plugins/my-plugin/status \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Documentation Phase

**Required Documentation**:

1. **README.md** - Main documentation
   - Quick description
   - Features list
   - Installation instructions
   - Configuration guide
   - Troubleshooting
   - Screenshots

2. **CHANGELOG.md** - Version history
   ```markdown
   # Changelog

   ## [1.0.0] - 2024-12-01
   ### Added
   - Initial release
   - Feature X
   - Feature Y
   ```

3. **API Documentation** (if applicable)
   - Endpoint descriptions
   - Request/response examples
   - Error codes

### 6. Release Phase

**Create Release Archive**:
```bash
# From plugin directory
tar czf ../my-plugin-v1.0.0.tar.gz \
  --exclude=".git" \
  --exclude="node_modules" \
  --exclude="*.log" \
  --exclude="tmp" \
  .
```

**GitHub Release**:
```bash
# Tag version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions will:
# 1. Create release
# 2. Upload tar.gz
# 3. Update registry.json
```

---

## ✅ Best Practices

### Security

**DO**:
- ✅ Use environment variables for secrets
- ✅ Validate all user input
- ✅ Run as non-root user (Docker)
- ✅ Use HTTPS for external connections
- ✅ Implement rate limiting
- ✅ Log security events

**DON'T**:
- ❌ Hardcode passwords/API keys
- ❌ Run as root unless absolutely necessary
- ❌ Trust user input without validation
- ❌ Expose internal services publicly
- ❌ Store secrets in logs

### Error Handling

```go
// Good error handling
result, err := doSomething()
if err != nil {
    log.Error().Err(err).Msg("Failed to do something")
    return fmt.Errorf("operation failed: %w", err)
}

// With context
if err := validateInput(input); err != nil {
    return errors.BadRequest("Invalid input", err)
}
```

### Logging

```go
// Use structured logging (zerolog)
log.Info().
    Str("plugin_id", pluginID).
    Int("port", port).
    Msg("Starting plugin")

log.Error().
    Err(err).
    Str("operation", "database_query").
    Msg("Database error")
```

### Configuration

```go
// Load from environment with defaults
func loadConfig() *Config {
    return &Config{
        Port:     getEnv("PORT", "8080"),
        LogLevel: getEnv("LOG_LEVEL", "info"),
        DBPath:   getEnv("DB_PATH", "/data/plugin.db"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### Database

```go
// Use GORM for SQLite
db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
if err != nil {
    log.Fatal().Err(err).Msg("Failed to open database")
}

// Auto-migrate models
db.AutoMigrate(&MyModel{})

// Queries
var results []MyModel
db.Where("status = ?", "active").Find(&results)
```

---

## 🎨 Code Style Guidelines

### Go

```go
// Package naming: lowercase, no underscores
package myplugin

// Exported functions: PascalCase
func ProcessRequest() {}

// Private functions: camelCase
func helperFunction() {}

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const DefaultTimeout = 30
const MAX_RETRIES = 3

// Interfaces: end with -er if possible
type Processor interface {
    Process() error
}

// Error handling: wrap errors with context
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}
```

### TypeScript/React

```typescript
// Components: PascalCase
export const PluginDashboard: React.FC = () => { }

// Hooks: start with "use"
const usePluginData = () => { }

// Types/Interfaces: PascalCase
interface PluginConfig {
  name: string;
  enabled: boolean;
}

// Constants: SCREAMING_SNAKE_CASE
const API_BASE_URL = "http://localhost:8080";
```

---

## 🧪 Testing Checklist

Before submitting, test:

- [ ] Plugin installs successfully
- [ ] Plugin starts without errors
- [ ] Health check endpoint responds
- [ ] All features work as expected
- [ ] Configuration can be changed
- [ ] Plugin can be stopped gracefully
- [ ] Plugin can be restarted
- [ ] Plugin can be uninstalled cleanly
- [ ] No port conflicts with other plugins
- [ ] Logs are clear and informative
- [ ] Error messages are helpful
- [ ] Documentation is complete and accurate

---

## 📦 Plugin Categories

Choose the most appropriate category:

- **storage** - Backup, sync, object storage (MinIO, Syncthing)
- **media** - Media servers, streaming (Plex, Jellyfin)
- **communication** - VoIP, chat, email (Asterisk, Matrix)
- **development** - Git, CI/CD, databases (Gitea, Jenkins)
- **monitoring** - Metrics, logs, alerts (Prometheus, Grafana)
- **networking** - DNS, VPN, proxy (Pi-hole, WireGuard)
- **productivity** - Tasks, notes, documents (Nextcloud, Bitwarden)
- **security** - Password managers, 2FA, firewalls
- **utilities** - General tools and utilities

---

## 🔗 Important Links

### Documentation
- **Main Docs**: https://docs.stumpfworks.com
- **Plugin Guide**: https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/DEVELOPMENT.md
- **API Reference**: https://docs.stumpfworks.com/api

### Repositories
- **Main NAS**: https://github.com/Stumpf-works/stumpfworks-nas
- **Plugin Registry**: https://github.com/Stumpf-works/stumpfworks-nas-apps
- **Example Plugins**: https://github.com/Stumpf-works/stumpfworks-nas-apps/tree/main/plugins

### Community
- **Discussions**: https://github.com/Stumpf-works/stumpfworks-nas/discussions
- **Issues**: https://github.com/Stumpf-works/stumpfworks-nas-apps/issues
- **Discord**: https://discord.gg/stumpfworks

---

## 🤖 AI Assistant Instructions

When helping with StumpfWorks NAS plugin development:

1. **Always** check if similar functionality exists in the core system or other plugins
2. **Prefer** using System Library APIs over direct shell commands
3. **Follow** established patterns from existing plugins
4. **Validate** plugin.json schema and required fields
5. **Ensure** proper error handling and logging
6. **Write** comprehensive documentation
7. **Consider** security implications of all code
8. **Test** thoroughly before suggesting completion
9. **Suggest** improvements to code quality and structure
10. **Reference** this master prompt for architectural decisions

---

## 📊 Common Patterns Reference

### Accessing StumpfWorks NAS API from Plugin

```go
package main

import (
    "fmt"
    "net/http"
    "os"
)

func callNasAPI(endpoint string) ([]byte, error) {
    apiURL := os.Getenv("NAS_API_URL")
    token := os.Getenv("NAS_API_TOKEN")

    req, err := http.NewRequest("GET", apiURL + endpoint, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer " + token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Configuration File Management

```go
// Read plugin.json config section
type Config struct {
    Setting1 string `json:"setting1"`
    Setting2 int    `json:"setting2"`
}

func loadPluginConfig() (*Config, error) {
    data, err := os.ReadFile("plugin.json")
    if err != nil {
        return nil, err
    }

    var manifest struct {
        Config Config `json:"config"`
    }

    if err := json.Unmarshal(data, &manifest); err != nil {
        return nil, err
    }

    return &manifest.Config, nil
}
```

### Docker Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1
```

---

## ✨ Example Plugin Structure

```
my-plugin/
├── plugin.json              # Manifest
├── README.md                # User docs
├── CHANGELOG.md             # Version history
├── LICENSE                  # License
├── docker-compose.yml       # Docker setup
│
├── backend/
│   ├── main.go             # Entry point
│   ├── api/
│   │   ├── server.go       # API server
│   │   └── handlers.go     # API handlers
│   ├── service/
│   │   └── service.go      # Business logic
│   ├── models/
│   │   └── models.go       # Data models
│   ├── config/
│   │   └── config.go       # Configuration
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/               # Optional
│   ├── src/
│   ├── package.json
│   └── Dockerfile
│
├── config/                # Config templates
│   └── default.conf
│
├── scripts/               # Utility scripts
│   └── setup.sh
│
└── docs/                  # Additional docs
    ├── API.md
    └── ARCHITECTURE.md
```

---

**End of Master Prompt**

Use this document as the authoritative source for StumpfWorks NAS plugin development. When in doubt, refer back to these guidelines.

**Good luck building amazing plugins! 🚀**
