# Plugin Name

Short description of your plugin.

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

Install via StumpfWorks NAS Plugin Store or manually:

```bash
sudo cp -r . /var/lib/stumpfworks/plugins/plugin-name/
```

## Configuration

Edit plugin.json config section:

```json
{
  "config": {
    "setting1": "value1"
  }
}
```

## Usage

### Via UI

1. Open StumpfWorks NAS
2. Navigate to Plugins
3. Find "Plugin Name"
4. Click Enable

### Via API

```bash
curl -X POST http://nas-ip:8080/api/v1/plugins/plugin-name/enable \
  -H "Authorization: Bearer $TOKEN"
```

## Troubleshooting

### Issue 1

Solution...

### Issue 2

Solution...

## License

MIT License - See [LICENSE](LICENSE)

## Credits

Built with ❤️ by [Your Name](https://github.com/username)
