# Hello World Plugin

A simple example plugin that demonstrates the StumpfWorks NAS Plugin System.

## Features

- Demonstrates plugin lifecycle (start, stop, restart)
- Shows how to read configuration from `plugin.json`
- Handles graceful shutdown on SIGTERM/SIGINT
- Periodic message logging based on configurable interval

## Building

```bash
# Build the plugin binary
go build -o hello-world main.go
```

## Installation

1. Build the plugin (see above)
2. Copy the entire `hello-world/` directory to `/var/lib/stumpfworks/plugins/`
3. Use the StumpfWorks NAS API or UI to enable and start the plugin

### Via API

```bash
# Enable the plugin
curl -X POST http://localhost:8080/api/v1/plugins/com.stumpfworks.hello-world/enable

# Start the plugin
curl -X POST http://localhost:8080/api/v1/plugins/com.stumpfworks.hello-world/start

# Check status
curl http://localhost:8080/api/v1/plugins/com.stumpfworks.hello-world/status

# Stop the plugin
curl -X POST http://localhost:8080/api/v1/plugins/com.stumpfworks.hello-world/stop
```

## Configuration

Edit `plugin.json` to customize:

- `message`: The message to print periodically
- `interval`: How often to print the message (in seconds)

## Environment Variables

The plugin receives these environment variables from the runtime:

- `PLUGIN_ID`: Unique plugin identifier
- `PLUGIN_DIR`: Plugin installation directory
- `NAS_API_URL`: StumpfWorks NAS API endpoint

## Development

This plugin serves as a template for developing your own plugins. Key concepts:

1. **Manifest** (`plugin.json`): Defines plugin metadata and configuration
2. **Entry Point**: Executable that runs when plugin starts (specified in manifest)
3. **Graceful Shutdown**: Handle SIGTERM/SIGINT for clean shutdown
4. **Logging**: Use stdout/stderr - they're automatically captured by the runtime
5. **Configuration**: Load settings from `plugin.json` config section

## Next Steps

- See `docs/PLUGIN_SDK.md` for full plugin development documentation
- Check `examples/plugins/` for more advanced examples
- Join the community to share your plugins!
