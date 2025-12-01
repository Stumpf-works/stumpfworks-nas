.PHONY: all build clean test run install docker dev

# Variables
BINARY_NAME=vpn-server
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build ${LDFLAGS} -o bin/$(BINARY_NAME) cmd/vpn-server/main.go
	@echo "Built: bin/$(BINARY_NAME)"

# Build for Linux (useful for cross-compilation from Windows)
build-linux:
	@echo "Building for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/$(BINARY_NAME)-linux-amd64 cmd/vpn-server/main.go
	@echo "Built: bin/$(BINARY_NAME)-linux-amd64"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "Cleaned!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run with race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

# Run the server
run: build
	@echo "Starting $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME) -config vpn-server.yaml

# Development mode with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server..."
	@air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run ./...

# Install the binary to system
install: build-linux
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp bin/$(BINARY_NAME)-linux-amd64 /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# Create systemd service file
service:
	@echo "Creating systemd service file..."
	@echo "[Unit]" > vpn-server.service
	@echo "Description=Stumpfworks VPN Server" >> vpn-server.service
	@echo "After=network.target" >> vpn-server.service
	@echo "" >> vpn-server.service
	@echo "[Service]" >> vpn-server.service
	@echo "Type=simple" >> vpn-server.service
	@echo "User=root" >> vpn-server.service
	@echo "WorkingDirectory=/opt/stumpfworks-nas/vpn-server" >> vpn-server.service
	@echo "ExecStart=/usr/local/bin/$(BINARY_NAME) -config /etc/stumpfworks-nas/vpn-server.yaml" >> vpn-server.service
	@echo "Restart=on-failure" >> vpn-server.service
	@echo "RestartSec=5s" >> vpn-server.service
	@echo "" >> vpn-server.service
	@echo "[Install]" >> vpn-server.service
	@echo "WantedBy=multi-user.target" >> vpn-server.service
	@echo "Service file created: vpn-server.service"
	@echo "Install with: sudo cp vpn-server.service /etc/systemd/system/"
	@echo "Enable with: sudo systemctl enable vpn-server"
	@echo "Start with: sudo systemctl start vpn-server"

# Build Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t stumpfworks/vpn-server:$(VERSION) .
	@echo "Docker image built: stumpfworks/vpn-server:$(VERSION)"

# Run in Docker
docker-run:
	@echo "Running in Docker..."
	@docker run -d \
		--name vpn-server \
		--cap-add=NET_ADMIN \
		--device /dev/net/tun \
		-p 8080:8080 \
		-p 51820:51820/udp \
		-v $(PWD)/vpn-server.yaml:/etc/vpn-server.yaml \
		stumpfworks/vpn-server:$(VERSION)

# Generate TLS certificates for development
certs:
	@echo "Generating self-signed certificates..."
	@mkdir -p certs
	@openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -days 365 -nodes -subj "/CN=localhost"
	@echo "Certificates generated in certs/"

# Database migrations
migrate:
	@echo "Running database migrations..."
	@go run cmd/vpn-server/main.go -config vpn-server.yaml -migrate
	@echo "Migrations complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  make build        - Build the binary"
	@echo "  make build-linux  - Build for Linux"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make run          - Build and run the server"
	@echo "  make dev          - Run in development mode with auto-reload"
	@echo "  make deps         - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make install      - Install binary to system"
	@echo "  make service      - Create systemd service file"
	@echo "  make docker       - Build Docker image"
	@echo "  make docker-run   - Run in Docker"
	@echo "  make certs        - Generate development certificates"
	@echo "  make help         - Show this help"
