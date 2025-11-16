.PHONY: help install dev build release test clean docker-build docker-up docker-down lint format

# Default target
help:
	@echo "Stumpf.Works NAS - Development Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make install       - Install all dependencies (backend + frontend)"
	@echo "  make dev           - Run development servers (backend + frontend)"
	@echo "  make build         - Build for production"
	@echo "  make release       - Build release binaries for all platforms"
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
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Copying frontend files for embedding..."
	mkdir -p backend/embedfs
	rm -rf backend/embedfs/dist
	cp -r frontend/dist backend/embedfs/
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

# ISO builder (future)
iso:
	@echo "ISO builder will be available in Phase 7"

# Generate API docs
docs:
	@echo "Generating API documentation..."
	cd backend && swag init -g cmd/stumpfworks-server/main.go
	@echo "✓ API docs generated"
