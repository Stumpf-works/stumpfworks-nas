.PHONY: help install dev build test clean docker-build docker-up docker-down lint format

# Default target
help:
	@echo "Stumpf.Works NAS - Development Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make install       - Install all dependencies (backend + frontend)"
	@echo "  make dev           - Run development servers (backend + frontend)"
	@echo "  make build         - Build for production"
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
	@echo "Building backend..."
	cd backend && go build -o ../dist/stumpfworks-server cmd/stumpfworks-server/main.go
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "✓ Build complete. Artifacts in ./dist/"

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

# ISO builder (future)
iso:
	@echo "ISO builder will be available in Phase 7"

# Generate API docs
docs:
	@echo "Generating API documentation..."
	cd backend && swag init -g cmd/stumpfworks-server/main.go
	@echo "✓ API docs generated"
