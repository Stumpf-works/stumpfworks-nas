# StumpfWorks NAS - API Audit Tool

Comprehensive API testing and audit tool for StumpfWorks NAS Backend.

## Features

- ✅ **Automatic Endpoint Discovery** - Tests all known API endpoints
- ✅ **Backend Function Testing** - Verifies system library functions
- ✅ **Prometheus Metrics Check** - Validates monitoring metrics
- ✅ **Performance Metriken** - Measures response times
- ✅ **Security Tests** - Tests authentication and authorization
- ✅ **Detailed Reports** - JSON + Markdown output

## Installation

```bash
cd tools/api-audit
go build -o audit-tool .
```

## Usage

### Basic Audit

```bash
# Get JWT token first
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.access_token')

# Run full audit
./audit-tool --url http://localhost:8080 --token "$TOKEN"
```

### Specific Tests

```bash
# Only test endpoints
./audit-tool --url http://localhost:8080 --token "$TOKEN" --endpoints-only

# Only test backend functions
./audit-tool --url http://localhost:8080 --token "$TOKEN" --backend-only

# Only test Prometheus metrics
./audit-tool --url http://localhost:8080 --token "$TOKEN" --metrics-only

# Test specific category
./audit-tool --url http://localhost:8080 --token "$TOKEN" --category storage

# Include destructive tests (DELETE endpoints)
./audit-tool --url http://localhost:8080 --token "$TOKEN" --force-destructive
```

### Advanced Options

```bash
# Custom output directory
./audit-tool --url http://localhost:8080 --token "$TOKEN" --output ./my_reports

# JSON only
./audit-tool --url http://localhost:8080 --token "$TOKEN" --format json

# Markdown only
./audit-tool --url http://localhost:8080 --token "$TOKEN" --format md

# Verbose output
./audit-tool --url http://localhost:8080 --token "$TOKEN" --verbose

# Custom timeout
./audit-tool --url http://localhost:8080 --token "$TOKEN" --timeout 30s
```

## Output

The tool generates two reports:

### 1. audit_report.json
Machine-readable JSON report containing:
- All test results
- Performance metrics
- Issues found
- Backend function status

### 2. AUDIT_API_REPORT.md
Human-readable Markdown report with:
- Executive summary
- Endpoint overview (passed/failed/skipped)
- Backend function status
- Prometheus metrics check
- Performance analysis
- Critical issues and recommendations

## Exit Codes

- `0` - All tests passed
- `1` - Some tests failed

## Tested Categories

### Endpoints
- Health & System
- Authentication
- Users & Groups
- Storage (ZFS, SMART)
- Sharing (Samba, NFS)
- Network
- Docker
- Files
- Backup
- Scheduler
- Metrics
- Updates
- Alerts
- Audit Logs
- 2FA
- Plugins

### Backend Functions
- ZFS pools and snapshots
- SMART disk health
- Samba service status
- Network interfaces
- DNS configuration
- Firewall status
- User management
- Group management

### Prometheus Metrics
- CPU usage
- Memory usage
- Network traffic
- Disk I/O
- ZFS pool health
- SMART disk status
- Service status
- Share connections

## Examples

### CI/CD Integration

```bash
#!/bin/bash
# ci-audit.sh

# Start NAS server
make dev &
sleep 5

# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.access_token')

# Run audit
./tools/api-audit/audit-tool \
  --url http://localhost:8080 \
  --token "$TOKEN" \
  --output ./audit_results \
  --format both

# Check exit code
if [ $? -ne 0 ]; then
  echo "❌ API audit failed"
  exit 1
fi

echo "✅ API audit passed"
```

### Scheduled Monitoring

```bash
# Run audit every hour and save results
0 * * * * cd /path/to/nas && ./tools/api-audit/audit-tool \
  --url http://localhost:8080 \
  --token "$TOKEN" \
  --output /var/log/nas-audits/$(date +\%Y-\%m-\%d-\%H)
```

## Development

### Add New Endpoint

Edit `internal/discovery/endpoints.go`:

```go
{
    Method:       "GET",
    Path:         "/api/v1/mynew/endpoint",
    Description:  "My new endpoint",
    AuthRequired: true,
    Destructive:  false,
},
```

### Add Custom Test

Edit `internal/testing/endpoint.go` to add custom validation logic.

### Add New Metric Check

Edit `internal/backend/system.go` to add new Prometheus metric checks.

## Troubleshooting

### "Connection refused"
Make sure the NAS backend is running:
```bash
cd backend && make run
```

### "401 Unauthorized"
Your token might be expired. Get a new one:
```bash
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.access_token')
```

### "404 Not Found"
Some endpoints might not be implemented yet. Check AUDIT_API_REPORT.md for missing endpoints.

## Contributing

To add new tests or improve the tool:

1. Add endpoint definitions in `internal/discovery/endpoints.go`
2. Implement custom tests in `internal/testing/`
3. Update report generation in `internal/report/generator.go`
4. Build and test: `go build && ./audit-tool --verbose`

## License

Copyright © 2025 Stumpf.Works. All rights reserved.
