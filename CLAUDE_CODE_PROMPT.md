# StumpfWorks NAS - Debian Packaging & Auto-Deployment

Ich habe ein APT-Repository unter http://apt.stumpfworks.de eingerichtet.
StumpfWorks NAS läuft bereits als manuell installierter Service.
Jetzt müssen wir es in ein professionelles Debian-Paket umwandeln und automatisch deployen.

## 🎯 ZIEL

Nutzer sollen "apt install stumpfworks-nas" ausführen können und bekommen ein 
produktionsreifes, vollständig konfiguriertes NAS-System mit PostgreSQL und 
einem professionellen CLI-Management-Tool (stumpfctl).

## 🌐 SERVER-ZUGRIFF

- Host: root@46.4.25.15
- SSH-Key: ~/.ssh/id_ed25519 (bereits konfiguriert)
- APT Repo Path: /var/www/apt-repo/pool/main/
- Update Command: update-apt-repo
- Du kannst SSH-Befehle für Deployment nutzen

## 📦 AKTUELLER STAND

- Binary: /usr/local/bin/stumpfworks-nas (Go mit embedded React)
- Service: /etc/systemd/system/stumpfworks-nas.service
- ÄNDERUNG: Migration von SQLite zu PostgreSQL
- Config: /etc/stumpfworks-nas/config.yaml
- Tools: diagnose, check-password, set-password

## 🛠️ TECH-STACK

**Backend:**
- Go mit Chi Router, GORM ORM
- PostgreSQL (statt SQLite) - WICHTIGE ÄNDERUNG!
- JWT Auth, 2FA, bcrypt
- Embedded Frontend (React build im Binary)

**Frontend:**
- React + TypeScript + Vite
- Wird ins Go Binary embedded (build-time)

**Features:**
- User Management (Admin/User/Guest)
- Samba Integration
- System Metrics & Monitoring
- Alerts (Email/Webhook)
- Backup & Scheduler
- Docker Integration
- Plugin System

## 🗄️ POSTGRESQL SETUP

### 1. Code-Änderungen

- GORM Driver von sqlite zu postgres wechseln
- Connection String: "host=localhost user=stumpfworks password=XXX dbname=stumpfworks_nas port=5432 sslmode=disable"
- Import: "gorm.io/driver/postgres"

### 2. Database Setup (via postinst)

- PostgreSQL Database "stumpfworks_nas" erstellen
- PostgreSQL User "stumpfworks" mit Passwort erstellen
- Grant Permissions
- Run Migrations automatisch

### 3. Vorteile

- Keine "welche DB"-Probleme mehr
- Bessere Concurrency
- Produktions-ready
- Einfacheres Backup (pg_dump)
- Multi-User Access ohne File-Locking

## 🏗️ PAKET-STRUKTUR

### 1. Haupt-Paket: stumpfworks-nas

**Dependencies:**
- postgresql (>= 13)
- postgresql-client
- samba
- systemd

### 2. Dateien installieren

```
/usr/bin/stumpfworks-nas           # Main binary (Server)
/usr/bin/stumpfctl                 # CLI Management Tool ⭐ WICHTIG!
/usr/bin/stumpfworks-diagnose      # Diagnostic tool
/usr/bin/stumpfworks-check-password
/usr/bin/stumpfworks-set-password
/usr/bin/stumpfworks-dbsetup       # NEU: PostgreSQL setup helper

/etc/stumpfworks-nas/
├── config.yaml                    # Haupt-Config
├── config.yaml.example            # Template
└── .db-password                   # Generiertes DB-Passwort (chmod 600)

/var/lib/stumpfworks-nas/          # Data directory
├── backups/                       # DB Backups
└── plugins/                       # Plugin directory

/var/log/stumpfworks-nas/          # Logs

/usr/share/stumpfworks-nas/
├── doc/                           # Documentation
├── examples/                      # Config examples
└── sql/
    └── init.sql                   # Optional: SQL Init Scripts

/etc/systemd/system/
└── stumpfworks-nas.service        # systemd service
```

## 📟 STUMPFCTL - CLI MANAGEMENT TOOL (⭐ KERN-FEATURE)

stumpfctl ist DAS Haupt-Management-Tool für StumpfWorks NAS.
Es ist ein separates Go-Binary das systemctl wrapper + viele zusätzliche Features bietet.

**WICHTIG:** systemctl funktioniert weiterhin parallel, aber stumpfctl ist die 
user-friendly, feature-reiche Alternative mit StumpfWorks-spezifischen Funktionen.

### 📋 STUMPFCTL COMMANDS

**Service Management (wrapper für systemctl):**
```bash
stumpfctl start                    # Service starten
stumpfctl stop                     # Service stoppen
stumpfctl restart                  # Service neu starten
stumpfctl reload                   # Config neu laden
stumpfctl status                   # Status + Health + Uptime + Connections
stumpfctl enable                   # Auto-start aktivieren
stumpfctl disable                  # Auto-start deaktivieren
```

**Logs & Monitoring:**
```bash
stumpfctl logs                     # Logs anzeigen (letzte 50 Zeilen)
stumpfctl logs --follow            # Live logs (wie tail -f)
stumpfctl logs -n 100              # Letzte 100 Zeilen
stumpfctl logs --since "1h ago"    # Logs der letzten Stunde
stumpfctl health                   # Umfassender Health Check
stumpfctl metrics                  # Aktuelle System-Metriken
stumpfctl stats                    # Performance-Statistiken
```

**Version & Info:**
```bash
stumpfctl version                  # Version (Server + CLI)
stumpfctl info                     # System-Info (OS, CPU, RAM, Disk)
stumpfctl about                    # About StumpfWorks NAS
```

**User Management:**
```bash
stumpfctl user list                # Alle User auflisten (Tabelle)
stumpfctl user add <username>      # Neuen User hinzufügen (interaktiv)
stumpfctl user add <username> --admin  # Admin-User erstellen
stumpfctl user delete <username>   # User löschen (mit Bestätigung)
stumpfctl user info <username>     # User-Details anzeigen
stumpfctl user password <username> # Password zurücksetzen (interaktiv)
stumpfctl user disable <username>  # User deaktivieren
stumpfctl user enable <username>   # User aktivieren
```

**Backup & Restore:**
```bash
stumpfctl backup create            # Manuelles DB Backup erstellen
stumpfctl backup list              # Alle Backups auflisten
stumpfctl backup restore <file>    # Backup wiederherstellen
stumpfctl backup cleanup           # Alte Backups löschen (>30 Tage)
stumpfctl backup auto              # Auto-Backup Status & Config
```

**Config Management:**
```bash
stumpfctl config show              # Komplette Config anzeigen (censored passwords)
stumpfctl config get <key>         # Einzelnen Wert anzeigen
stumpfctl config set <key> <value> # Wert setzen
stumpfctl config edit              # Config in $EDITOR öffnen
stumpfctl config validate          # Config validieren
stumpfctl config reset             # Config auf Defaults zurücksetzen
```

**Share Management:**
```bash
stumpfctl share list               # Alle Shares auflisten
stumpfctl share add <name>         # Neuen Share erstellen (interaktiv)
stumpfctl share delete <name>      # Share löschen
stumpfctl share info <name>        # Share-Details anzeigen
stumpfctl share reload             # Samba Config neu laden
```

**System Maintenance:**
```bash
stumpfctl update check             # Nach Updates suchen
stumpfctl update install           # Updates installieren
stumpfctl doctor                   # System-Diagnose (umfassend)
stumpfctl reset-admin              # Admin-Password zurücksetzen (Recovery)
```

### 🎨 STUMPFCTL OUTPUT-DESIGN

**Beispiel-Output von "stumpfctl status":**

```
┌─────────────────────────────────────────────────────────┐
│              StumpfWorks NAS Status                     │
├─────────────────────────────────────────────────────────┤
│ Service:      ● Running (healthy)                       │
│ Uptime:       3 days, 5 hours, 23 minutes              │
│ Version:      v0.1.0 (build 2025-11-20)                │
│ PID:          1234                                      │
├─────────────────────────────────────────────────────────┤
│ Database:     ✓ Connected (PostgreSQL 15.3)            │
│ Web UI:       ✓ Accessible (http://localhost:8080)     │
│ Samba:        ✓ Running (3 active shares)              │
│ Active Users: 5 logged in                              │
├─────────────────────────────────────────────────────────┤
│ Health Score: 95/100 ✓                                  │
│ CPU:          12% (4 cores)                            │
│ Memory:       1.2 GB / 8 GB (15%)                      │
│ Disk:         456 GB / 2 TB (23%)                      │
└─────────────────────────────────────────────────────────┘

Verwende "stumpfctl logs" für Logs
Verwende "stumpfctl health" für detaillierte Diagnose
```

**Beispiel "stumpfctl user list":**

```
┌──────────┬────────┬─────────┬──────────────────────┬────────┐
│ Username │ Role   │ Status  │ Last Login           │ 2FA    │
├──────────┼────────┼─────────┼──────────────────────┼────────┤
│ admin    │ Admin  │ Active  │ 2025-11-20 10:30:15 │ ✓ Yes  │
│ john     │ User   │ Active  │ 2025-11-19 14:22:01 │ ✗ No   │
│ jane     │ User   │ Disabled│ 2025-11-18 09:15:42 │ ✗ No   │
│ guest    │ Guest  │ Active  │ Never               │ ✗ No   │
└──────────┴────────┴─────────┴──────────────────────┴────────┘

Total: 4 users (3 active, 1 disabled)
```

### 🔧 STUMPFCTL IMPLEMENTATION

**Code-Struktur:**
```
backend/
├── cmd/
│   ├── stumpfworks-server/        # Haupt-Service
│   └── stumpfctl/                 # CLI Tool ⭐
│       ├── main.go
│       └── commands/
│           ├── service.go         # start, stop, restart, status
│           ├── logs.go            # log commands
│           ├── user.go            # User Management
│           ├── backup.go          # Backup Commands
│           ├── config.go          # Config Management
│           ├── share.go           # Share Management
│           ├── health.go          # Health & Metrics
│           └── system.go          # System Info & Maintenance
└── pkg/
    ├── cli/                       # Shared CLI utilities
    │   ├── table.go              # Pretty tables
    │   ├── prompt.go             # Interactive prompts
    │   └── color.go              # Colored output
    └── client/                    # API Client für stumpfctl
        └── api.go                # HTTP Client zum Server
```

**Go Libraries für stumpfctl:**
- github.com/spf13/cobra         # CLI Framework
- github.com/fatih/color         # Colored Output
- github.com/olekukonko/tablewriter  # Pretty Tables
- github.com/manifoldco/promptui # Interactive Prompts
- github.com/briandowns/spinner  # Loading Spinners

**Features:**
- Ruft systemctl commands via exec.Command()
- Kommuniziert mit Server via HTTP API (localhost:8080)
- Liest Config aus /etc/stumpfworks-nas/config.yaml
- Colorized Output (✓ grün, ✗ rot, ● gelb)
- Interactive Prompts für gefährliche Operationen
- Progress Bars für lange Operationen
- Auto-Completion Support (bash/zsh)

## 📝 CONFIG (config.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  environment: "production"

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  database: "stumpfworks_nas"
  username: "stumpfworks"
  password: "${DB_PASSWORD}"  # Aus /etc/stumpfworks-nas/.db-password
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

logging:
  path: "/var/log/stumpfworks-nas"
  level: "info"
  format: "json"  # oder "text"

samba:
  enabled: true
  config_path: "/etc/samba/smb.conf"
  shares_path: "/etc/stumpfworks-nas/shares.conf"

backup:
  enabled: true
  schedule: "0 2 * * *"  # Daily at 2 AM
  retention_days: 30
  path: "/var/lib/stumpfworks-nas/backups"

alerts:
  email:
    enabled: false
    smtp_host: ""
    smtp_port: 587
  webhook:
    enabled: false
    url: ""
```

## 🚀 DEPLOYMENT-STRATEGIE (3-stufig)

### A) MAKEFILE TARGETS

```makefile
.PHONY: all clean build frontend backend tools deb deploy release test-install

# Version aus Git
VERSION := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

all: clean build

clean:
	rm -rf dist/
	rm -rf frontend/dist/

frontend:
	cd frontend && npm install && npm run build

backend: frontend
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="$(LDFLAGS)" \
		-o dist/stumpfworks-nas \
		./backend/cmd/stumpfworks-server

tools:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/stumpfctl ./backend/cmd/stumpfctl
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/stumpfworks-diagnose ./backend/cmd/diagnose
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/stumpfworks-check-password ./backend/cmd/check-password
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/stumpfworks-set-password ./backend/cmd/set-password
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/stumpfworks-dbsetup ./backend/cmd/dbsetup

completions:
	./dist/stumpfctl completion bash > dist/stumpfctl-completion.bash
	./dist/stumpfctl completion zsh > dist/stumpfctl-completion.zsh

build: backend tools completions

deb: build
	./scripts/build-deb.sh $(VERSION)

deploy: deb
	./scripts/deploy.sh $(VERSION)

release: deploy
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

test-install:
	sudo apt install ./dist/stumpfworks-nas_$(VERSION)_amd64.deb

test-remove:
	sudo apt remove stumpfworks-nas

test-purge:
	sudo apt purge stumpfworks-nas
```

### B) DEPLOY SCRIPT (scripts/deploy.sh)

```bash
#!/bin/bash
set -e

VERSION=${1:-$(git describe --tags --always)}
DEB_FILE="dist/stumpfworks-nas_${VERSION}_amd64.deb"
SERVER="root@46.4.25.15"
REPO_PATH="/var/www/apt-repo/pool/main/"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 StumpfWorks NAS Deployment v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if .deb exists
if [ ! -f "$DEB_FILE" ]; then
    echo "❌ Error: $DEB_FILE not found!"
    echo "Run 'make deb' first."
    exit 1
fi

echo "📤 Uploading $DEB_FILE to APT server..."
scp "$DEB_FILE" "$SERVER:$REPO_PATH"

echo "🔄 Updating repository metadata..."
ssh "$SERVER" "update-apt-repo"

echo "✅ Deployment complete!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Package available at:"
echo "   http://apt.stumpfworks.de"
echo ""
echo "📥 Install with:"
echo "   sudo apt update"
echo "   sudo apt install stumpfworks-nas"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Verify package in repository
echo ""
echo "🔍 Verifying package in repository..."
if ssh "$SERVER" "apt-cache policy stumpfworks-nas" | grep -q "$VERSION"; then
    echo "✓ Package verified successfully!"
else
    echo "⚠️  Warning: Could not verify package version"
fi
```

### C) GITHUB ACTION (.github/workflows/release.yml)

```yaml
name: Build and Deploy Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for git describe
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
      
      - name: Get version
        id: version
        run: |
          VERSION=$(git describe --tags --always | sed 's/^v//')
          echo "version=$VERSION" >> $GITHUB_OUTPUT
          echo "Building version: $VERSION"
      
      - name: Build Frontend
        run: |
          cd frontend
          npm install
          npm run build
      
      - name: Build Backend
        run: |
          make build
      
      - name: Build Debian Package
        run: |
          make deb VERSION=${{ steps.version.outputs.version }}
      
      - name: Setup SSH
        env:
          SSH_PRIVATE_KEY: ${{ secrets.APT_SERVER_SSH_KEY }}
        run: |
          mkdir -p ~/.ssh
          echo "$SSH_PRIVATE_KEY" > ~/.ssh/id_ed25519
          chmod 600 ~/.ssh/id_ed25519
          ssh-keyscan -H 46.4.25.15 >> ~/.ssh/known_hosts
      
      - name: Deploy to APT Server
        run: |
          make deploy VERSION=${{ steps.version.outputs.version }}
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*.deb
          body: |
            ## StumpfWorks NAS ${{ steps.version.outputs.version }}
            
            ### Installation
            ```bash
            sudo apt update
            sudo apt install stumpfworks-nas
            ```
            
            ### What's New
            See [CHANGELOG.md](CHANGELOG.md) for details.
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## ⚡ AUFGABEN - SCHRITT FÜR SCHRITT

### PHASE 1 - Code Migration
1. Migriere Code von SQLite zu PostgreSQL
   - Ändere GORM Driver
   - Teste Migrations
   - Passe Connection-Handling an
   
2. Erstelle stumpfworks-dbsetup Tool:
   - DB-Setup automatisieren
   - Password generieren
   - Credentials testen

### PHASE 2 - stumpfctl Development
3. Erstelle cmd/stumpfctl/ mit allen Commands
4. Implementiere pkg/cli/ utilities
5. Implementiere pkg/client/ API client
6. Teste alle stumpfctl Commands

### PHASE 3 - Packaging
7. Vervollständige debian/ Verzeichnis
8. Schreibe postinst/prerm/postrm Scripts mit PostgreSQL-Setup
9. Erstelle systemd service file
10. Fixe config.yaml Handling (Development vs Production paths)

### PHASE 4 - Build System
11. Erstelle/Erweitere Makefile mit allen Targets
12. Erstelle deploy.sh Script
13. Erstelle build-deb.sh Script
14. Baue erste .deb (Version 0.1.0)
15. Teste Deployment auf APT Server

### PHASE 5 - Automation
16. Erstelle GitHub Action für automatisches Release-Building
17. Teste kompletten Workflow: Git Tag → Auto Build → Auto Deploy

### PHASE 6 - Testing
18. Teste Installation auf frischem Debian System
19. Verifiziere PostgreSQL Setup
20. Teste Admin-Login
21. Teste alle CLI Tools (stumpfctl)
22. Teste Deinstallation/Purge

## 🎯 ERSTE VERSION SOLL

- ✅ PostgreSQL automatisch einrichten
- ✅ Cleanly installieren
- ✅ DB Migrations automatisch ausführen
- ✅ Service auto-starten
- ✅ Admin-Login funktionieren
- ✅ Alle Tools verfügbar sein (inkl. stumpfctl)
- ✅ Sauber deinstallierbar sein
- ✅ DB-Backups bei Deinstallation anbieten
- ✅ Via "make deploy" auf APT Server deploybar
- ✅ GitHub Actions für Auto-Release

## 💡 DEPLOYMENT-FLOWS

**Entwickler (manuell):**
```bash
git tag v0.1.0
make release        # Build + Deploy alles
```

**CI/CD (automatisch):**
```bash
git tag v0.1.0
git push --tags     # GitHub Action übernimmt Rest
```

**Endnutzer:**
```bash
sudo apt update
sudo apt install stumpfworks-nas
sudo systemctl status stumpfworks-nas
stumpfctl status
```

## 🔍 VERIFICATION

**Nach Deployment prüfen:**
```bash
curl http://apt.stumpfworks.de/dists/stable/main/binary-amd64/Packages | grep stumpfworks-nas
apt-cache policy stumpfworks-nas
```

**Nach Installation prüfen:**
```bash
systemctl status stumpfworks-nas
stumpfctl version
stumpfctl status
curl http://localhost:8080/health
```

## 🚦 START MIT

1. Zeige mir geplante Projekt-Struktur (debian/, scripts/, .github/)
2. PostgreSQL Migration Code-Changes
3. Dann bauen wir Schritt für Schritt alle Phasen
4. Am Ende: Komplettes Auto-Deployment System

**LOS GEHT'S! 🚀**
