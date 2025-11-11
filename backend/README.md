# Stumpf.Works NAS - Backend

Go-based backend server for Stumpf.Works NAS Solution.

## Features

✅ **Phase 2 Complete** - Backend Core Infrastructure

- RESTful API with Chi router
- JWT-based authentication
- User management (CRUD operations)
- System information API (CPU, RAM, Disk, Network)
- WebSocket support for real-time events
- SQLite database with GORM ORM
- Structured logging with Zap
- Configuration management with Viper
- Graceful shutdown
- systemd integration

## Project Structure

```
backend/
├── cmd/
│   └── stumpfworks-server/    # Main application entry point
├── internal/                   # Private application code
│   ├── api/                   # HTTP API layer
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # HTTP middleware (auth, logging)
│   │   └── websocket/         # WebSocket handling
│   ├── config/                # Configuration management
│   ├── database/              # Database layer
│   │   └── models/            # Database models
│   ├── system/                # System information
│   └── users/                 # User management
└── pkg/                       # Public libraries
    ├── errors/                # Error handling
    ├── logger/                # Logging infrastructure
    └── utils/                 # Utility functions
```

## Getting Started

### Prerequisites

- Go 1.21+
- SQLite3

### Installation

```bash
# Clone repository (if not already done)
cd backend

# Download dependencies
go mod download

# Copy example config (optional)
cp config.example.yaml config.yaml
```

### Running

```bash
# Development mode
go run cmd/stumpfworks-server/main.go

# Or use Make (from project root)
make dev-backend
```

The server will start on `http://localhost:8080`.

### Building

```bash
# Build binary
go build -o stumpfworks-server cmd/stumpfworks-server/main.go

# Or use Make (from project root)
make build
```

## Configuration

Configuration can be provided via:
1. **Config file**: `config.yaml` (or path set in `STUMPFWORKS_CONFIG` env var)
2. **Environment variables**: Prefix with `STUMPFWORKS_`

Example `config.yaml`:

```yaml
app:
  name: "Stumpf.Works NAS"
  version: "0.1.0-alpha"
  environment: "development"

server:
  host: "0.0.0.0"
  port: 8080

database:
  driver: "sqlite"
  path: "./data/stumpfworks.db"

auth:
  jwtSecret: "${STUMPFWORKS_AUTH_JWTSECRET}"
  jwtExpirationHours: 24

logging:
  level: "info"
  development: true
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/` | API information |
| POST | `/api/v1/auth/login` | User login |

### Protected Endpoints (Require Authentication)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/logout` | User logout | User |
| POST | `/api/v1/auth/refresh` | Refresh token | User |
| GET | `/api/v1/auth/me` | Get current user | User |
| GET | `/api/v1/system/info` | System information | User |
| GET | `/api/v1/system/metrics` | Real-time metrics | User |
| GET | `/api/v1/users` | List users | Admin |
| POST | `/api/v1/users` | Create user | Admin |
| GET | `/api/v1/users/{id}` | Get user | Admin |
| PUT | `/api/v1/users/{id}` | Update user | Admin |
| DELETE | `/api/v1/users/{id}` | Delete user | Admin |

### WebSocket

| Endpoint | Description |
|----------|-------------|
| `/ws` | WebSocket connection |

## Authentication

The API uses JWT (JSON Web Tokens) for authentication.

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin"
  }'
```

Response:

```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@stumpfworks.local",
      "role": "admin",
      "isActive": true
    }
  }
}
```

### Using the Token

Include the access token in the `Authorization` header:

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Default Credentials

**⚠️ Security Warning:** Change these immediately in production!

- **Username:** `admin`
- **Password:** `admin`

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/users
```

## Linting & Formatting

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

## Deployment

### Systemd Service

1. Build binary:
   ```bash
   go build -o stumpfworks-server cmd/stumpfworks-server/main.go
   ```

2. Copy to installation directory:
   ```bash
   sudo cp stumpfworks-server /opt/stumpfworks/
   sudo chmod +x /opt/stumpfworks/stumpfworks-server
   ```

3. Install systemd service:
   ```bash
   sudo cp ../systemd/stumpfworks.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable stumpfworks
   sudo systemctl start stumpfworks
   ```

4. Check status:
   ```bash
   sudo systemctl status stumpfworks
   ```

## Development

### Adding New Endpoints

1. Create handler in `internal/api/handlers/`
2. Add route in `internal/api/router.go`
3. Implement business logic in appropriate package (e.g., `internal/users/`)

### Adding Database Models

1. Create model in `internal/database/models/`
2. Add to `RunMigrations()` in `internal/database/migrations.go`

## Troubleshooting

### Database locked error

SQLite has limited concurrent write support. If you get "database is locked" errors:
- Ensure only one instance is running
- Consider switching to PostgreSQL for production

### Permission denied errors

Ensure the application has write permissions to:
- Database directory (`./data/`)
- Log directory (if file logging is enabled)

## Next Steps (Phase 3)

- Frontend development (React + TailwindCSS)
- Storage management endpoints
- Network management endpoints
- Plugin system implementation

## License

MIT License - see [LICENSE](../LICENSE) for details.
