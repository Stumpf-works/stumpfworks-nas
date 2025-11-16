# StumpfWorks NAS Plugin SDK

Complete guide for developing plugins for StumpfWorks NAS v1.1.0+

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Plugin Structure](#plugin-structure)
- [Plugin Manifest](#plugin-manifest)
- [Plugin Lifecycle](#plugin-lifecycle)
- [Environment Variables](#environment-variables)
- [Configuration](#configuration)
- [Logging](#logging)
- [Best Practices](#best-practices)
- [API Reference](#api-reference)
- [Examples](#examples)

---

## Overview

The StumpfWorks NAS Plugin System allows developers to extend the NAS functionality with custom applications. Plugins run as separate processes, providing isolation and security.

### Key Features

- **Sandboxed Execution**: Plugins run in separate processes
- **Lifecycle Management**: Start, stop, restart plugins via API
- **Configuration**: JSON-based configuration system
- **Logging**: Automatic stdout/stderr capture
- **Environment**: Access to system information via environment variables
- **Language Agnostic**: Write plugins in any language (Go, Python, Node.js, etc.)

---

## Getting Started

### Prerequisites

- StumpfWorks NAS v1.1.0 or higher
- Development environment (Go, Python, Node.js, etc.)
- Basic understanding of REST APIs

### Quick Start

1. **Create Plugin Directory**:
   ```bash
   mkdir my-plugin
   cd my-plugin
   ```

2. **Create Manifest** (`plugin.json`):
   ```json
   {
     "id": "com.example.my-plugin",
     "name": "My Plugin",
     "version": "1.0.0",
     "author": "Your Name",
     "description": "My awesome plugin",
     "entryPoint": "my-plugin",
     "config": {}
   }
   ```

3. **Create Entry Point** (executable):
   - Go: Compile to binary
   - Python: Use shebang `#!/usr/bin/env python3`
   - Node.js: Use shebang `#!/usr/bin/env node`

4. **Make Executable**:
   ```bash
   chmod +x my-plugin
   ```

5. **Install Plugin**:
   ```bash
   cp -r my-plugin /var/lib/stumpfworks/plugins/
   ```

---

## Plugin Structure

```
my-plugin/
├── plugin.json       # Manifest (required)
├── my-plugin         # Entry point executable (required)
├── README.md         # Documentation (recommended)
├── config/           # Additional config files (optional)
├── assets/           # Images, icons, etc. (optional)
└── lib/              # Libraries, dependencies (optional)
```

---

## Plugin Manifest

The `plugin.json` file defines your plugin's metadata and configuration.

### Schema

```json
{
  "id": "string",              // Unique identifier (required)
  "name": "string",            // Display name (required)
  "version": "string",         // Semantic version (required)
  "author": "string",          // Author name (required)
  "description": "string",     // Short description (required)
  "icon": "string",            // Icon (emoji or path) (optional)
  "entryPoint": "string",      // Executable filename (required)
  "config": {}                 // Default configuration (optional)
}
```

### ID Convention

Use reverse domain notation: `com.company.plugin-name`

**Examples**:
- `com.stumpfworks.backup-scheduler`
- `org.example.monitoring-agent`
- `io.github.username.custom-app`

### Version

Follow [Semantic Versioning](https://semver.org/):
- `1.0.0` - Major.Minor.Patch
- `0.1.0-beta` - Pre-release versions

---

## Plugin Lifecycle

### States

1. **Installed**: Plugin files present, not enabled
2. **Enabled**: Plugin enabled, but not running
3. **Running**: Plugin process active
4. **Stopped**: Plugin was running, now stopped
5. **Crashed**: Plugin exited with error

### Lifecycle Events

```
Install → Enable → Start → Running
                      ↓
                   Stop ← Restart
                      ↓
                  Stopped
```

### API Endpoints

```bash
# Enable plugin
POST /api/v1/plugins/{id}/enable

# Start plugin
POST /api/v1/plugins/{id}/start

# Stop plugin
POST /api/v1/plugins/{id}/stop

# Restart plugin
POST /api/v1/plugins/{id}/restart

# Get status
GET /api/v1/plugins/{id}/status
```

---

## Environment Variables

The runtime provides these environment variables to your plugin:

| Variable | Description | Example |
|----------|-------------|---------|
| `PLUGIN_ID` | Unique plugin identifier | `com.example.my-plugin` |
| `PLUGIN_DIR` | Plugin installation directory | `/var/lib/stumpfworks/plugins/com.example.my-plugin` |
| `NAS_API_URL` | StumpfWorks NAS API endpoint | `http://localhost:8080/api/v1` |

### Usage Example (Go)

```go
pluginID := os.Getenv("PLUGIN_ID")
pluginDir := os.Getenv("PLUGIN_DIR")
apiURL := os.Getenv("NAS_API_URL")
```

### Usage Example (Python)

```python
import os

plugin_id = os.getenv("PLUGIN_ID")
plugin_dir = os.getenv("PLUGIN_DIR")
api_url = os.getenv("NAS_API_URL")
```

---

## Configuration

### Loading Configuration

Configuration is stored in the `config` section of `plugin.json`:

```json
{
  "config": {
    "interval": 60,
    "enabled_features": ["feature1", "feature2"],
    "api_key": ""
  }
}
```

### Example (Go)

```go
type Config struct {
    Interval        int      `json:"interval"`
    EnabledFeatures []string `json:"enabled_features"`
    APIKey          string   `json:"api_key"`
}

func loadConfig() (*Config, error) {
    pluginDir := os.Getenv("PLUGIN_DIR")
    data, _ := os.ReadFile(filepath.Join(pluginDir, "plugin.json"))

    var manifest struct {
        Config Config `json:"config"`
    }
    json.Unmarshal(data, &manifest)
    return &manifest.Config, nil
}
```

### Updating Configuration

Users can update configuration via API:

```bash
curl -X PUT http://localhost:8080/api/v1/plugins/{id}/config \
  -H "Content-Type: application/json" \
  -d '{"config": {"interval": 120}}'
```

**Note**: Restart plugin for config changes to take effect.

---

## Logging

### Standard Streams

Plugin output is automatically captured:

- **stdout** → Logged at INFO level
- **stderr** → Logged at ERROR level

### Example (Go)

```go
fmt.Println("This goes to INFO logs")
fmt.Fprintln(os.Stderr, "This goes to ERROR logs")
```

### Example (Python)

```python
print("This goes to INFO logs")
print("This goes to ERROR logs", file=sys.stderr)
```

### Best Practices

1. Use structured logging (JSON) for machine parsing
2. Include timestamps in log messages
3. Log important events (start, stop, errors, warnings)
4. Avoid excessive logging (rate limit if needed)

---

## Best Practices

### Security

1. **Validate All Input**: Never trust user configuration
2. **Least Privilege**: Request minimal permissions
3. **Secrets Management**: Don't store secrets in plain text
4. **Rate Limiting**: Limit API calls and resource usage

### Performance

1. **Graceful Shutdown**: Handle SIGTERM/SIGINT properly
2. **Resource Limits**: Respect memory and CPU constraints
3. **Error Handling**: Don't crash on errors, log and recover
4. **Async Operations**: Use goroutines/async for I/O

### Code Quality

1. **Error Handling**: Always check and handle errors
2. **Documentation**: Include README and code comments
3. **Testing**: Write unit tests for critical functions
4. **Versioning**: Follow semantic versioning

### Example: Graceful Shutdown (Go)

```go
func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Main loop
    for {
        select {
        case <-sigChan:
            cleanup()
            return
        default:
            doWork()
        }
    }
}
```

---

## API Reference

### Plugin Management

#### List Plugins

```
GET /api/v1/plugins
```

Response:
```json
[
  {
    "id": "com.example.my-plugin",
    "name": "My Plugin",
    "version": "1.0.0",
    "enabled": true,
    "installed": true
  }
]
```

#### Get Plugin

```
GET /api/v1/plugins/{id}
```

#### Install Plugin

```
POST /api/v1/plugins/install
Content-Type: application/json

{
  "sourcePath": "/path/to/plugin"
}
```

#### Enable/Disable Plugin

```
POST /api/v1/plugins/{id}/enable
POST /api/v1/plugins/{id}/disable
```

#### Update Configuration

```
PUT /api/v1/plugins/{id}/config
Content-Type: application/json

{
  "config": {
    "key": "value"
  }
}
```

### Plugin Runtime

#### Start Plugin

```
POST /api/v1/plugins/{id}/start
```

#### Stop Plugin

```
POST /api/v1/plugins/{id}/stop
```

#### Restart Plugin

```
POST /api/v1/plugins/{id}/restart
```

#### Get Status

```
GET /api/v1/plugins/{id}/status
```

Response:
```json
{
  "pluginID": "com.example.my-plugin",
  "status": "running",
  "running": true,
  "startedAt": "2025-11-16T10:00:00Z"
}
```

#### List Running Plugins

```
GET /api/v1/plugins/running
```

---

## Examples

### Example 1: Hello World (Go)

See `examples/plugins/hello-world/` for a complete example.

### Example 2: System Monitor (Python)

```python
#!/usr/bin/env python3
import os
import time
import psutil
import signal
import sys

def handle_signal(signum, frame):
    print("Shutting down gracefully...")
    sys.exit(0)

signal.signal(signal.SIGINT, handle_signal)
signal.signal(signal.SIGTERM, handle_signal)

print(f"Starting plugin: {os.getenv('PLUGIN_ID')}")

while True:
    cpu = psutil.cpu_percent()
    mem = psutil.virtual_memory().percent
    print(f"CPU: {cpu}%, Memory: {mem}%")
    time.sleep(5)
```

### Example 3: Notification Service (Node.js)

```javascript
#!/usr/bin/env node
const fs = require('fs');
const path = require('path');

const pluginID = process.env.PLUGIN_ID;
const pluginDir = process.env.PLUGIN_DIR;

console.log(`Starting plugin: ${pluginID}`);

// Load configuration
const manifest = JSON.parse(
  fs.readFileSync(path.join(pluginDir, 'plugin.json'), 'utf8')
);

// Handle shutdown
process.on('SIGTERM', () => {
  console.log('Received SIGTERM, shutting down');
  process.exit(0);
});

// Main loop
setInterval(() => {
  console.log('Plugin is running...');
}, manifest.config.interval * 1000);
```

---

## Troubleshooting

### Plugin Won't Start

1. Check executable permissions: `chmod +x plugin-binary`
2. Verify entry point in `plugin.json` matches filename
3. Check logs: `journalctl -f | grep "Plugin:"`

### Plugin Crashes Immediately

1. Test plugin manually: `./plugin-binary`
2. Check for missing dependencies
3. Review stderr logs for errors

### Configuration Not Loading

1. Verify JSON syntax in `plugin.json`
2. Ensure config section exists
3. Check file permissions

---

## Support

- **Documentation**: https://docs.stumpf.works
- **Examples**: `examples/plugins/` directory
- **Issues**: https://github.com/Stumpf-works/stumpfworks-nas/issues
- **Community**: https://github.com/Stumpf-works/stumpfworks-nas/discussions

---

## License

Plugin SDK is part of StumpfWorks NAS and is licensed under the MIT License.

---

**Happy Plugin Development! 🚀**
