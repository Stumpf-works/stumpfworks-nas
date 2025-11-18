.PHONY: help install dev build release test clean docker-build docker-up docker-down lint format upgrade install-system uninstall

# Default target
help:
	@echo "Stumpf.Works NAS - Development Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make install       - Install all dependencies (backend + frontend)"
	@echo "  make dev           - Run development servers (backend + frontend)"
	@echo "  make build         - Build for production"
	@echo "  make release       - Build release binaries for all platforms"
	@echo ""
	@echo "System Installation:"
	@echo "  make install-system - Install to system (creates systemd service)"
	@echo "  make upgrade        - Upgrade existing installation (auto-detects location)"
	@echo "  make uninstall      - Remove system installation"
	@echo ""
	@echo "  make test          - Run all tests"
	@echo "  make lint          - Run linters"
	@echo "  make format        - Format code"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start Docker Compose stack"
	@echo "  make docker-down   - Stop Docker Compose stack"

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✓ Dependencies installed"

# Run development servers
dev:
	@echo "Starting development servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@make -j2 dev-backend dev-frontend

dev-backend:
	cd backend && go run cmd/stumpfworks-server/main.go

dev-frontend:
	cd frontend && npm run dev

# Build for production
build:
	@echo "Checking frontend dependencies..."
	@if [ ! -d "frontend/node_modules" ]; then \
		echo "Installing frontend dependencies..."; \
		cd frontend && npm install; \
	fi
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Copying frontend files for embedding..."
	mkdir -p backend/embedfs
	rm -rf backend/embedfs/dist
	cp -r frontend/dist backend/embedfs/
	@echo "Checking backend dependencies..."
	@cd backend && go mod tidy
	@cd backend && go mod download
	@echo "Building backend with embedded frontend..."
	mkdir -p dist
	cd backend && go build -ldflags="-s -w" -o ../dist/stumpfworks-server cmd/stumpfworks-server/main.go
	@echo "✓ Build complete. Binary in ./dist/stumpfworks-server (includes embedded frontend)"

# Run tests
test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend && npm run test
	@echo "✓ All tests passed"

# Run linters
lint:
	@echo "Linting backend..."
	cd backend && golangci-lint run
	@echo "Linting frontend..."
	cd frontend && npm run lint
	@echo "✓ Linting complete"

# Format code
format:
	@echo "Formatting backend..."
	cd backend && go fmt ./...
	@echo "Formatting frontend..."
	cd frontend && npm run format
	@echo "✓ Formatting complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/
	rm -rf backend/dist/
	rm -rf frontend/dist/
	@echo "✓ Clean complete"

# Docker targets
docker-build:
	@echo "Building Docker images..."
	docker-compose build
	@echo "✓ Docker images built"

docker-up:
	@echo "Starting Docker Compose stack..."
	docker-compose up -d
	@echo "✓ Stack running"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

docker-down:
	@echo "Stopping Docker Compose stack..."
	docker-compose down
	@echo "✓ Stack stopped"

# Build release binaries for multiple platforms
release:
	@echo "Building frontend for release..."
	cd frontend && npm run build
	@echo "Copying frontend files for embedding..."
	mkdir -p backend/embedfs
	rm -rf backend/embedfs/dist
	cp -r frontend/dist backend/embedfs/
	@echo "Checking backend dependencies..."
	@cd backend && go mod tidy
	@cd backend && go mod download
	@echo "Building release binaries for multiple platforms..."
	@mkdir -p dist/releases

	# Linux AMD64
	@echo "Building for Linux AMD64..."
	cd backend && GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -X main.AppVersion=$(shell git describe --tags --always)" \
		-o ../dist/releases/stumpfworks-nas-linux-amd64 \
		cmd/stumpfworks-server/main.go

	# Linux ARM64
	@echo "Building for Linux ARM64..."
	cd backend && GOOS=linux GOARCH=arm64 go build \
		-ldflags="-s -w -X main.AppVersion=$(shell git describe --tags --always)" \
		-o ../dist/releases/stumpfworks-nas-linux-arm64 \
		cmd/stumpfworks-server/main.go

	# Linux ARM (Raspberry Pi)
	@echo "Building for Linux ARM..."
	cd backend && GOOS=linux GOARCH=arm GOARM=7 go build \
		-ldflags="-s -w -X main.AppVersion=$(shell git describe --tags --always)" \
		-o ../dist/releases/stumpfworks-nas-linux-armv7 \
		cmd/stumpfworks-server/main.go

	# Darwin AMD64 (Intel Mac)
	@echo "Building for Darwin AMD64..."
	cd backend && GOOS=darwin GOARCH=amd64 go build \
		-ldflags="-s -w -X main.AppVersion=$(shell git describe --tags --always)" \
		-o ../dist/releases/stumpfworks-nas-darwin-amd64 \
		cmd/stumpfworks-server/main.go

	# Darwin ARM64 (Apple Silicon)
	@echo "Building for Darwin ARM64..."
	cd backend && GOOS=darwin GOARCH=arm64 go build \
		-ldflags="-s -w -X main.AppVersion=$(shell git describe --tags --always)" \
		-o ../dist/releases/stumpfworks-nas-darwin-arm64 \
		cmd/stumpfworks-server/main.go

	# Build frontend
	@echo "Building frontend..."
	cd frontend && npm run build
	@cp -r frontend/dist dist/releases/frontend

	# Create checksums
	@echo "Creating checksums..."
	cd dist/releases && sha256sum stumpfworks-nas-* > checksums.txt

	@echo "✓ Release build complete. Binaries in ./dist/releases/"
	@ls -lh dist/releases/

# Upgrade existing installation
upgrade:
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║  Stumpf.Works NAS - Intelligent Upgrade System                ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "🔍 Detecting installation..."
	@# Detect architecture
	@ARCH=$$(uname -m); \
	case $$ARCH in \
		x86_64) BINARY_NAME="stumpfworks-nas-linux-amd64" ;; \
		aarch64) BINARY_NAME="stumpfworks-nas-linux-arm64" ;; \
		armv7l) BINARY_NAME="stumpfworks-nas-linux-armv7" ;; \
		*) echo "❌ Unsupported architecture: $$ARCH"; exit 1 ;; \
	esac; \
	echo "   Architecture: $$ARCH ($$BINARY_NAME)"; \
	\
	INSTALL_PATH=""; \
	SERVICE_NAME=""; \
	SYSTEMD_SERVICE=""; \
	\
	if [ -f /usr/local/bin/stumpfworks-nas ]; then \
		INSTALL_PATH="/usr/local/bin/stumpfworks-nas"; \
		echo "   Found binary: /usr/local/bin/stumpfworks-nas"; \
	elif [ -f /usr/bin/stumpfworks-nas ]; then \
		INSTALL_PATH="/usr/bin/stumpfworks-nas"; \
		echo "   Found binary: /usr/bin/stumpfworks-nas"; \
	elif [ -f /opt/stumpfworks-nas/stumpfworks-nas ]; then \
		INSTALL_PATH="/opt/stumpfworks-nas/stumpfworks-nas"; \
		echo "   Found binary: /opt/stumpfworks-nas/stumpfworks-nas"; \
	else \
		echo "❌ No existing installation found!"; \
		echo "   Searched: /usr/local/bin, /usr/bin, /opt/stumpfworks-nas"; \
		echo "   Run 'make install-system' for first-time installation"; \
		exit 1; \
	fi; \
	\
	if systemctl list-units --full --all | grep -q "stumpfworks-nas.service"; then \
		SYSTEMD_SERVICE="stumpfworks-nas.service"; \
		echo "   Found systemd service: stumpfworks-nas.service"; \
	elif systemctl list-units --full --all | grep -q "stumpfworks.service"; then \
		SYSTEMD_SERVICE="stumpfworks.service"; \
		echo "   Found systemd service: stumpfworks.service"; \
	fi; \
	\
	echo ""; \
	echo "📦 Building new version..."; \
	$(MAKE) build GOOS=linux GOARCH=$$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/;s/armv7l/arm/') 2>&1 | grep -E "(Building|✓|Copying)" || true; \
	\
	if [ ! -f "dist/stumpfworks-server" ]; then \
		echo "❌ Build failed!"; \
		exit 1; \
	fi; \
	\
	echo ""; \
	echo "💾 Creating backup..."; \
	BACKUP_PATH="$$INSTALL_PATH.backup.$$(date +%Y%m%d-%H%M%S)"; \
	cp "$$INSTALL_PATH" "$$BACKUP_PATH" 2>/dev/null || sudo cp "$$INSTALL_PATH" "$$BACKUP_PATH"; \
	echo "   Backup saved: $$BACKUP_PATH"; \
	\
	if [ -n "$$SYSTEMD_SERVICE" ]; then \
		echo ""; \
		echo "⏸️  Stopping service..."; \
		sudo systemctl stop "$$SYSTEMD_SERVICE" 2>/dev/null || true; \
		sleep 2; \
		if systemctl is-active --quiet "$$SYSTEMD_SERVICE"; then \
			echo "   ⚠️  Service still running, waiting..."; \
			sleep 3; \
		fi; \
		echo "   ✓ Service stopped"; \
	fi; \
	\
	echo ""; \
	echo "📥 Installing new version..."; \
	cp dist/stumpfworks-server "$$INSTALL_PATH" 2>/dev/null || sudo cp dist/stumpfworks-server "$$INSTALL_PATH"; \
	chmod +x "$$INSTALL_PATH" 2>/dev/null || sudo chmod +x "$$INSTALL_PATH"; \
	echo "   ✓ Binary updated"; \
	\
	NEW_VERSION=$$(dist/stumpfworks-server --version 2>/dev/null | head -1 || echo "v1.3.0"); \
	echo "   New version: $$NEW_VERSION"; \
	\
	if [ -n "$$SYSTEMD_SERVICE" ]; then \
		echo ""; \
		echo "▶️  Starting service..."; \
		sudo systemctl start "$$SYSTEMD_SERVICE"; \
		sleep 2; \
		if systemctl is-active --quiet "$$SYSTEMD_SERVICE"; then \
			echo "   ✓ Service started successfully"; \
			STATUS=$$(systemctl status "$$SYSTEMD_SERVICE" --no-pager -l | head -3 | tail -1); \
			echo "   $$STATUS"; \
		else \
			echo "   ❌ Service failed to start!"; \
			echo "   Rolling back..."; \
			sudo cp "$$BACKUP_PATH" "$$INSTALL_PATH"; \
			sudo systemctl start "$$SYSTEMD_SERVICE"; \
			echo "   ✓ Rollback complete"; \
			exit 1; \
		fi; \
	fi; \
	\
	echo ""; \
	echo "╔════════════════════════════════════════════════════════════════╗"; \
	echo "║  ✅ Upgrade Complete!                                          ║"; \
	echo "╚════════════════════════════════════════════════════════════════╝"; \
	echo ""; \
	echo "📍 Installation: $$INSTALL_PATH"; \
	echo "💾 Backup: $$BACKUP_PATH"; \
	if [ -n "$$SYSTEMD_SERVICE" ]; then \
		echo "🔄 Service: $$SYSTEMD_SERVICE (running)"; \
	fi; \
	echo ""; \
	echo "🌐 Access your NAS:"; \
	echo "   http://localhost:8080"; \
	echo "   http://$$(hostname -I | awk '{print $$1}'):8080"; \
	echo "";

# Install to system (first-time installation)
install-system:
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║  Stumpf.Works NAS - System Installation                       ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@# Check if already installed
	@if [ -f /usr/local/bin/stumpfworks-nas ] || [ -f /usr/bin/stumpfworks-nas ] || [ -f /opt/stumpfworks-nas/stumpfworks-nas ]; then \
		echo "⚠️  Installation detected!"; \
		echo ""; \
		echo "Found existing installation. Use 'make upgrade' instead."; \
		echo "Or run 'make uninstall' first to remove the old installation."; \
		exit 1; \
	fi
	@echo "📦 Building Stumpf.Works NAS..."
	@$(MAKE) build 2>&1 | grep -E "(Building|✓|Copying)" || true
	@if [ ! -f "dist/stumpfworks-server" ]; then \
		echo "❌ Build failed!"; \
		exit 1; \
	fi
	@echo ""
	@echo "📥 Installing to system..."
	@# Create installation directory
	@sudo mkdir -p /opt/stumpfworks-nas
	@sudo mkdir -p /etc/stumpfworks-nas
	@sudo mkdir -p /var/lib/stumpfworks-nas
	@sudo mkdir -p /var/log/stumpfworks-nas
	@# Install binary
	@sudo cp dist/stumpfworks-server /usr/local/bin/stumpfworks-nas
	@sudo chmod +x /usr/local/bin/stumpfworks-nas
	@echo "   ✓ Binary installed to /usr/local/bin/stumpfworks-nas"
	@# Create default config if it doesn't exist
	@if [ ! -f /etc/stumpfworks-nas/config.yaml ]; then \
		echo "   Creating default configuration..."; \
		JWT_SECRET=$$(openssl rand -base64 32) && \
		printf '%s\n' \
			'# Stumpf.Works NAS Configuration' \
			'server:' \
			'  host: 0.0.0.0' \
			'  port: 8080' \
			'  environment: production' \
			'  allowedOrigins: []' \
			'' \
			'database:' \
			'  path: /var/lib/stumpfworks-nas/stumpfworks.db' \
			'' \
			'storage:' \
			'  basePath: /mnt/storage' \
			'  shares:' \
			'    - name: "Public"' \
			'      path: "/mnt/storage/public"' \
			'      readOnly: false' \
			'' \
			'auth:' \
			"  jwtSecret: \"$$JWT_SECRET\"" \
			'  sessionTimeout: 24h' \
			'' \
			'logging:' \
			'  level: info' \
			'  file: /var/log/stumpfworks-nas/stumpfworks.log' \
		| sudo tee /etc/stumpfworks-nas/config.yaml > /dev/null && \
		echo "   ✓ Configuration created at /etc/stumpfworks-nas/config.yaml"; \
	else \
		echo "   ⚠️  Existing config found, keeping it"; \
	fi
	@# Create systemd service
	@echo "   Creating systemd service..."
	@printf '%s\n' \
		'[Unit]' \
		'Description=Stumpf.Works NAS Server' \
		'After=network.target' \
		'' \
		'[Service]' \
		'Type=simple' \
		'User=root' \
		'Group=root' \
		'WorkingDirectory=/opt/stumpfworks-nas' \
		'ExecStart=/usr/local/bin/stumpfworks-nas --config /etc/stumpfworks-nas/config.yaml' \
		'Restart=always' \
		'RestartSec=10' \
		'StandardOutput=journal' \
		'StandardError=journal' \
		'SyslogIdentifier=stumpfworks-nas' \
		'' \
		'# Security settings' \
		'NoNewPrivileges=false' \
		'PrivateTmp=true' \
		'' \
		'# Resource limits' \
		'LimitNOFILE=65536' \
		'' \
		'[Install]' \
		'WantedBy=multi-user.target' \
	| sudo tee /etc/systemd/system/stumpfworks-nas.service > /dev/null
	@echo "   ✓ Systemd service created"
	@# Reload systemd
	@sudo systemctl daemon-reload
	@echo "   ✓ Systemd daemon reloaded"
	@# Enable and start service
	@echo ""
	@echo "🚀 Starting service..."
	@sudo systemctl enable stumpfworks-nas.service
	@sudo systemctl start stumpfworks-nas.service
	@sleep 2
	@if systemctl is-active --quiet stumpfworks-nas.service; then \
		echo "   ✓ Service started successfully"; \
	else \
		echo "   ❌ Service failed to start!"; \
		echo "   Check logs with: sudo journalctl -u stumpfworks-nas.service"; \
		exit 1; \
	fi
	@echo ""
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║  ✅ Installation Complete!                                     ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📍 Installed to: /usr/local/bin/stumpfworks-nas"
	@echo "⚙️  Configuration: /etc/stumpfworks-nas/config.yaml"
	@echo "💾 Database: /var/lib/stumpfworks-nas/stumpfworks.db"
	@echo "📄 Logs: /var/log/stumpfworks-nas/stumpfworks.log"
	@echo "🔄 Service: stumpfworks-nas.service (enabled & running)"
	@echo ""
	@echo "🌐 Access your NAS:"
	@echo "   http://localhost:8080"
	@echo "   http://$$(hostname -I | awk '{print $$1}'):8080"
	@echo ""
	@echo "📋 Useful commands:"
	@echo "   sudo systemctl status stumpfworks-nas    - Check status"
	@echo "   sudo systemctl restart stumpfworks-nas   - Restart service"
	@echo "   sudo journalctl -u stumpfworks-nas -f    - View logs"
	@echo "   make upgrade                              - Upgrade to new version"
	@echo "   make uninstall                            - Remove installation"
	@echo ""
	@echo "⚠️  IMPORTANT: Change the default admin password after first login!"
	@echo ""

# Uninstall from system
uninstall:
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║  Stumpf.Works NAS - Uninstallation                            ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "⚠️  This will remove Stumpf.Works NAS from your system."
	@echo "⚠️  Your data and configuration will be preserved."
	@echo ""
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" != "yes" ]; then \
		echo "Aborted."; \
		exit 0; \
	fi
	@echo ""
	@echo "🛑 Stopping service..."
	@sudo systemctl stop stumpfworks-nas.service 2>/dev/null || true
	@sudo systemctl disable stumpfworks-nas.service 2>/dev/null || true
	@echo "   ✓ Service stopped"
	@echo ""
	@echo "🗑️  Removing files..."
	@sudo rm -f /usr/local/bin/stumpfworks-nas
	@sudo rm -f /usr/bin/stumpfworks-nas
	@sudo rm -f /etc/systemd/system/stumpfworks-nas.service
	@sudo systemctl daemon-reload
	@echo "   ✓ Binary removed"
	@echo "   ✓ Systemd service removed"
	@echo ""
	@echo "📁 Preserved files (remove manually if needed):"
	@if [ -d /etc/stumpfworks-nas ]; then \
		echo "   /etc/stumpfworks-nas/ (configuration)"; \
	fi
	@if [ -d /var/lib/stumpfworks-nas ]; then \
		echo "   /var/lib/stumpfworks-nas/ (database)"; \
	fi
	@if [ -d /var/log/stumpfworks-nas ]; then \
		echo "   /var/log/stumpfworks-nas/ (logs)"; \
	fi
	@if [ -d /opt/stumpfworks-nas ]; then \
		echo "   /opt/stumpfworks-nas/ (data)"; \
	fi
	@echo ""
	@echo "✅ Uninstallation complete!"
	@echo ""
	@echo "To completely remove all data:"
	@echo "  sudo rm -rf /etc/stumpfworks-nas"
	@echo "  sudo rm -rf /var/lib/stumpfworks-nas"
	@echo "  sudo rm -rf /var/log/stumpfworks-nas"
	@echo "  sudo rm -rf /opt/stumpfworks-nas"
	@echo ""

# ISO builder (future)
iso:
	@echo "ISO builder will be available in Phase 7"

# Generate API docs
docs:
	@echo "Generating API documentation..."
	cd backend && swag init -g cmd/stumpfworks-server/main.go
	@echo "✓ API docs generated"
