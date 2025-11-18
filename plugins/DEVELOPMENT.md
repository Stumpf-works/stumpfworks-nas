# Plugin Development Guide

Detaillierte Anleitung zur Entwicklung von StumpfWorks NAS Plugins.

## 📚 Inhaltsverzeichnis

- [Einführung](#einführung)
- [Entwicklungsumgebung](#entwicklungsumgebung)
- [Plugin-Typen](#plugin-typen)
- [Schritt-für-Schritt Tutorial](#schritt-für-schritt-tutorial)
- [Best Practices](#best-practices)
- [API Integration](#api-integration)
- [Testing](#testing)
- [Deployment](#deployment)

## 🎯 Einführung

StumpfWorks NAS Plugins sind eigenständige Anwendungen, die:
- Als separate Prozesse laufen
- Über REST APIs mit StumpfWorks kommunizieren
- Optional Docker-Container verwenden können
- Eine eigene UI bereitstellen können

## 🛠️ Entwicklungsumgebung

### Voraussetzungen

- **Go 1.21+** (für Backend)
- **Node.js 18+** (für Frontend, optional)
- **Docker & Docker Compose** (für Container-basierte Plugins)
- **StumpfWorks NAS Development Setup**

### Entwicklungs-Workflow

```bash
# 1. StumpfWorks NAS klonen
git clone https://github.com/stumpf-works/stumpfworks-nas.git
cd stumpfworks-nas

# 2. Plugin-Ordner erstellen
mkdir -p plugins/my-plugin
cd plugins/my-plugin

# 3. Plugin entwickeln
# ... siehe Tutorial unten

# 4. Plugin lokal testen
./scripts/install-plugin.sh my-plugin

# 5. StumpfWorks NAS neu starten
systemctl restart stumpfworks-nas
```

## 🔧 Plugin-Typen

### 1. Standalone Go Plugin

**Wann verwenden:** Einfache Plugins ohne externe Dependencies

**Struktur:**
```
my-plugin/
├── plugin.json
├── main.go
└── go.mod
```

**Beispiel:**
```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    pluginID := os.Getenv("PLUGIN_ID")
    log.Printf("Starting plugin: %s", pluginID)

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 2. Docker-basiertes Plugin

**Wann verwenden:** Plugins mit komplexen Dependencies (z.B. Asterisk, PostgreSQL)

**Struktur:**
```
my-plugin/
├── plugin.json
├── docker-compose.yml
├── backend/
│   ├── main.go
│   └── Dockerfile
└── config/
    └── service.conf
```

### 3. Hybrid Plugin (Go + Frontend)

**Wann verwenden:** Plugins mit eigener Web-UI

**Struktur:**
```
my-plugin/
├── plugin.json
├── backend/
│   ├── main.go
│   └── api/
├── frontend/
│   ├── package.json
│   ├── src/
│   └── dist/
└── docker-compose.yml
```

## 📝 Schritt-für-Schritt Tutorial

### Beispiel: "System Monitor" Plugin

Wir erstellen ein Plugin, das Systemmetriken sammelt und über eine API bereitstellt.

#### Schritt 1: Plugin-Struktur erstellen

```bash
mkdir -p plugins/system-monitor/{backend,config}
cd plugins/system-monitor
```

#### Schritt 2: plugin.json erstellen

```json
{
  "id": "com.stumpfworks.system-monitor",
  "name": "System Monitor",
  "version": "1.0.0",
  "author": "StumpfWorks Team",
  "description": "Advanced system monitoring and metrics",
  "icon": "📊",
  "entryPoint": "system-monitor",
  "requires": {
    "docker": false,
    "ports": [9100],
    "storage": "100MB",
    "minNasVersion": "0.1.0"
  },
  "config": {
    "interval": 60,
    "metrics": ["cpu", "memory", "disk", "network"]
  }
}
```

#### Schritt 3: Go Backend entwickeln

**backend/main.go:**
```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
    Timestamp   time.Time `json:"timestamp"`
    CPUPercent  float64   `json:"cpu_percent"`
    MemoryUsed  uint64    `json:"memory_used"`
    MemoryTotal uint64    `json:"memory_total"`
}

var latestMetrics Metrics

func collectMetrics() {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

    for {
        cpuPercent, _ := cpu.Percent(time.Second, false)
        memInfo, _ := mem.VirtualMemory()

        latestMetrics = Metrics{
            Timestamp:   time.Now(),
            CPUPercent:  cpuPercent[0],
            MemoryUsed:  memInfo.Used,
            MemoryTotal: memInfo.Total,
        }

        <-ticker.C
    }
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(latestMetrics)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    pluginID := os.Getenv("PLUGIN_ID")
    log.Printf("Starting %s", pluginID)

    // Start metrics collection
    go collectMetrics()

    // Setup HTTP routes
    r := mux.NewRouter()
    r.HandleFunc("/health", healthCheck).Methods("GET")
    r.HandleFunc("/metrics", getMetrics).Methods("GET")

    port := ":9100"
    log.Printf("Listening on %s", port)
    log.Fatal(http.ListenAndServe(port, r))
}
```

**backend/go.mod:**
```go
module system-monitor

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/shirou/gopsutil/v3 v3.24.1
)
```

#### Schritt 4: Build-Skript erstellen

**build.sh:**
```bash
#!/bin/bash
set -e

echo "Building System Monitor Plugin..."

cd backend
go mod download
go build -o ../system-monitor main.go

echo "✅ Build complete! Executable: ./system-monitor"
```

```bash
chmod +x build.sh
```

#### Schritt 5: Plugin bauen und testen

```bash
# Build
./build.sh

# Lokal testen
PLUGIN_ID=com.stumpfworks.system-monitor ./system-monitor

# In anderem Terminal testen:
curl http://localhost:9100/health
curl http://localhost:9100/metrics
```

#### Schritt 6: Installation

```bash
# Plugin installieren
sudo mkdir -p /var/lib/stumpfworks/plugins/system-monitor
sudo cp -r . /var/lib/stumpfworks/plugins/system-monitor/

# Plugin über StumpfWorks API aktivieren
curl -X POST http://localhost:8080/api/v1/plugins/system-monitor/enable \
  -H "Authorization: Bearer $TOKEN"
```

## 🎨 Frontend Integration

### Plugin mit eigener UI

Wenn dein Plugin eine UI benötigt, kannst du eine React-App erstellen:

#### Frontend-Struktur

```
frontend/
├── package.json
├── vite.config.ts
├── src/
│   ├── App.tsx
│   ├── api/
│   │   └── client.ts
│   └── components/
│       └── Dashboard.tsx
└── dist/ (nach build)
```

#### Frontend in StumpfWorks registrieren

**Option 1: Iframe Embedding**
```json
{
  "ui": {
    "type": "iframe",
    "url": "http://localhost:9100/ui"
  }
}
```

**Option 2: Native Integration**
```typescript
// In StumpfWorks frontend/src/apps/index.tsx registrieren:
import { SystemMonitor } from '@/plugins/system-monitor/SystemMonitor';

export const registeredApps: App[] = [
  // ... existing apps
  {
    id: 'system-monitor',
    name: 'System Monitor',
    icon: '📊',
    component: SystemMonitor,
    defaultSize: { width: 1200, height: 800 },
  },
];
```

## 🔐 API Integration

### StumpfWorks API nutzen

Plugins können auf die StumpfWorks API zugreifen:

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func getUsers() ([]byte, error) {
    apiURL := os.Getenv("NAS_API_URL")
    token := os.Getenv("NAS_API_TOKEN")

    req, _ := http.NewRequest("GET", apiURL+"/users", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Verfügbare APIs

- **Users**: `/api/v1/users`
- **Storage**: `/api/v1/storage/*`
- **Docker**: `/api/v1/docker/*`
- **System Library**: `/api/v1/syslib/*`

## 🐳 Docker Integration

### Docker Compose verwenden

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  plugin-service:
    image: my-plugin:latest
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "9100:9100"
    environment:
      - PLUGIN_ID=${PLUGIN_ID}
      - NAS_API_URL=${NAS_API_URL}
    volumes:
      - plugin-data:/data
    restart: unless-stopped

volumes:
  plugin-data:
```

### StumpfWorks startet Docker Compose automatisch

Wenn `docker-compose.yml` vorhanden ist, startet StumpfWorks das Plugin automatisch via Docker Compose.

## ✅ Best Practices

### 1. Fehlerbehandlung

```go
// ❌ Schlecht
data, _ := fetchData()

// ✅ Gut
data, err := fetchData()
if err != nil {
    log.Printf("Error fetching data: %v", err)
    return err
}
```

### 2. Konfiguration

```go
// Konfiguration aus plugin.json lesen
type Config struct {
    Interval int      `json:"interval"`
    Metrics  []string `json:"metrics"`
}

func loadConfig() (*Config, error) {
    data, err := os.ReadFile("plugin.json")
    if err != nil {
        return nil, err
    }

    var manifest struct {
        Config Config `json:"config"`
    }

    json.Unmarshal(data, &manifest)
    return &manifest.Config, nil
}
```

### 3. Logging

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("Plugin started", "plugin_id", os.Getenv("PLUGIN_ID"))
```

### 4. Health Checks

Jedes Plugin sollte einen `/health` Endpoint haben:

```go
func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "version": "1.0.0",
    })
}
```

### 5. Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":9100"}

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)

    log.Println("Plugin stopped gracefully")
}
```

## 🧪 Testing

### Unit Tests

```go
// backend/metrics_test.go
package main

import "testing"

func TestCollectMetrics(t *testing.T) {
    metrics := collectMetrics()

    if metrics.CPUPercent < 0 || metrics.CPUPercent > 100 {
        t.Errorf("Invalid CPU percentage: %f", metrics.CPUPercent)
    }
}
```

### Integration Tests

```bash
# test.sh
#!/bin/bash

# Start plugin
./system-monitor &
PID=$!

sleep 2

# Test health endpoint
STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9100/health)
if [ "$STATUS" != "200" ]; then
    echo "❌ Health check failed"
    kill $PID
    exit 1
fi

echo "✅ All tests passed"
kill $PID
```

## 📦 Deployment

### 1. Release erstellen

```bash
# Version in plugin.json erhöhen
# Build erstellen
./build.sh

# Release-Archiv erstellen
tar czf system-monitor-v1.0.0.tar.gz \
    plugin.json \
    system-monitor \
    README.md \
    LICENSE
```

### 2. Installation dokumentieren

**README.md:**
```markdown
## Installation

### Automatisch (empfohlen)
curl -sSL https://plugins.stumpfworks.com/install.sh | bash -s system-monitor

### Manuell
1. Download: https://github.com/.../system-monitor-v1.0.0.tar.gz
2. Entpacken: tar xzf system-monitor-v1.0.0.tar.gz
3. Installieren: sudo cp -r system-monitor /var/lib/stumpfworks/plugins/
4. Aktivieren über StumpfWorks UI
```

## 🔍 Debugging

### Logs anzeigen

```bash
# Plugin-Logs
journalctl -u stumpfworks-nas -f | grep system-monitor

# Docker Logs (bei Container-Plugins)
docker logs -f system-monitor
```

### Common Issues

**Problem:** Plugin startet nicht
```bash
# Prüfe Berechtigungen
ls -la /var/lib/stumpfworks/plugins/my-plugin/

# Prüfe Executable
file /var/lib/stumpfworks/plugins/my-plugin/my-plugin
```

**Problem:** API-Zugriff fehlgeschlagen
```bash
# Prüfe Umgebungsvariablen
echo $NAS_API_URL
echo $NAS_API_TOKEN
```

## 📚 Weitere Ressourcen

- [Asterisk VoIP Plugin](./asterisk-voip/) - Vollständiges Beispiel
- [StumpfWorks API Docs](../docs/API.md)
- [Plugin SDK Reference](../docs/PLUGIN_SDK.md)

## 🤝 Support

- GitHub Issues: https://github.com/stumpf-works/stumpfworks-nas/issues
- Discussions: https://github.com/stumpf-works/stumpfworks-nas/discussions
- Discord: https://discord.gg/stumpfworks
