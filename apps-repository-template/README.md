# StumpfWorks NAS - Official Plugin Repository 🎉

Welcome to the official plugin repository for StumpfWorks NAS! This repository contains vetted, community-contributed plugins that extend StumpfWorks NAS with additional functionality.

[![Plugins](https://img.shields.io/badge/plugins-1+-blue)](./plugins/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

---

## 📦 Available Plugins

### Communication
- 📞 **[Asterisk VoIP](./plugins/asterisk-voip/)** - Complete VoIP telephone system with PBX features

### Storage (Coming Soon)
- ☁️ **MinIO** - S3-compatible object storage
- 🔄 **Syncthing** - Continuous file synchronization

### Media (Coming Soon)
- 🎬 **Plex** - Media server
- 📺 **Jellyfin** - Open-source media server

### Development (Coming Soon)
- 🦊 **Gitea** - Self-hosted Git service
- 🔨 **Jenkins** - CI/CD automation

---

## 🚀 Quick Start

### For Users

#### Install from StumpfWorks NAS UI

1. Open StumpfWorks NAS web interface
2. Navigate to **Plugins** → **Store**
3. Browse available plugins
4. Click **Install** on desired plugin
5. Configure and enable

#### Install via CLI

```bash
# List available plugins
curl http://your-nas:8080/api/v1/store/plugins

# Install a plugin
curl -X POST http://your-nas:8080/api/v1/store/plugins/com.stumpfworks.asterisk-voip/install \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🛠️ For Plugin Developers

Want to create a plugin? We'd love your contribution!

### 1. Read the Docs

- [Plugin Development Guide](https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/DEVELOPMENT.md)
- [Contributing Guidelines](./CONTRIBUTING.md)
- [Plugin Templates](./templates/)

### 2. Use a Template

```bash
# Clone this repository
git clone https://github.com/Stumpf-works/stumpfworks-nas-apps.git
cd stumpfworks-nas-apps

# Copy a template
cp -r templates/basic-plugin plugins/my-plugin
cd plugins/my-plugin

# Customize plugin.json
nano plugin.json
```

### 3. Develop Your Plugin

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed instructions.

### 4. Submit for Review

1. Fork this repository
2. Add your plugin to `plugins/your-plugin/`
3. Create a Pull Request
4. Pass automated checks
5. Await review from maintainers

---

## 📋 Plugin Categories

- **🔐 Storage & Backup** - Object storage, sync, backup solutions
- **🎬 Media** - Media servers, streaming, photo management
- **📞 Communication** - VoIP, chat, email, collaboration
- **💻 Development** - Git, CI/CD, databases, IDEs
- **📊 Monitoring** - Metrics, logging, alerting
- **🌐 Networking** - DNS, VPN, reverse proxy, ad-blocking
- **✅ Productivity** - Task management, notes, documents
- **🔒 Security** - Password managers, 2FA, firewalls
- **🛠️ Utilities** - Various tools and utilities

---

## 🔍 How to Find Plugins

### Browse by Category

```bash
curl "http://your-nas:8080/api/v1/store/plugins/search?category=communication"
```

### Search by Keyword

```bash
curl "http://your-nas:8080/api/v1/store/plugins/search?q=voip"
```

### View in UI

The StumpfWorks NAS web interface provides a beautiful plugin store with:
- 🖼️ Screenshots
- ⭐ Ratings & reviews
- 📊 Download statistics
- 🔄 One-click installation & updates

---

## 📊 Registry System

This repository uses a centralized `registry.json` file that StumpfWorks NAS instances query to discover plugins.

### Registry URL

```
https://raw.githubusercontent.com/Stumpf-works/stumpfworks-nas-apps/main/registry.json
```

### Registry Format

```json
{
  "version": "1.0.0",
  "updated": "2024-12-01T10:00:00Z",
  "plugins": [
    {
      "id": "com.stumpfworks.plugin-name",
      "name": "Plugin Name",
      "version": "1.0.0",
      "download_url": "https://...",
      ...
    }
  ]
}
```

The registry is automatically updated via GitHub Actions when:
- A new release is tagged
- A plugin is added/updated
- Daily (scheduled sync)

---

## 🤖 Automated Checks

All plugins undergo automated validation:

- ✅ **Syntax Check** - `plugin.json` validation
- 🔍 **Static Analysis** - Code quality checks
- 🐳 **Docker Validation** - docker-compose.yml syntax
- 🔒 **Security Scan** - Malware and vulnerability detection
- 📝 **Documentation Check** - README, CHANGELOG presence
- 🧪 **Basic Tests** - Plugin can be installed/started

---

## 📈 Statistics

| Metric | Count |
|--------|-------|
| Total Plugins | 1 |
| Categories | 9 |
| Total Downloads | 0 |
| Contributors | 1 |

*Updated automatically by GitHub Actions*

---

## 🤝 Contributing

We welcome contributions from everyone! Here's how to get started:

1. **Ideas** - Open an [issue](https://github.com/Stumpf-works/stumpfworks-nas-apps/issues) to discuss your plugin idea
2. **Develop** - Follow our [development guide](./CONTRIBUTING.md)
3. **Test** - Ensure your plugin works on StumpfWorks NAS
4. **Submit** - Create a pull request
5. **Review** - Work with maintainers to refine your plugin
6. **Publish** - Once approved, your plugin goes live!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

---

## 🔒 Security

### Reporting Security Issues

If you discover a security vulnerability in a plugin, please email:
**security@stumpfworks.com**

Do NOT open a public issue for security vulnerabilities.

### Plugin Review Process

All plugins are reviewed for:
- Malicious code
- Security vulnerabilities
- Resource abuse
- Data privacy concerns
- Compliance with guidelines

Approved plugins receive a ✅ badge.

---

## 📄 License

Individual plugins have their own licenses (see each plugin's LICENSE file).

This repository structure and documentation: [MIT License](LICENSE)

---

## 🙏 Credits

StumpfWorks NAS Apps is maintained by:
- **StumpfWorks Team** - Core maintainers
- **Community Contributors** - Plugin developers

Special thanks to all contributors who make this ecosystem possible!

---

## 📚 Resources

- **Main Repository**: https://github.com/Stumpf-works/stumpfworks-nas
- **Documentation**: https://docs.stumpfworks.com
- **Community Forum**: https://community.stumpfworks.com
- **Discord**: https://discord.gg/stumpfworks

---

## ❓ FAQ

**Q: How do I install plugins?**
A: Use the Plugin Store in StumpfWorks NAS UI or the REST API.

**Q: Are plugins sandboxed?**
A: Docker-based plugins run in containers. Native plugins run with limited permissions.

**Q: Can I monetize my plugin?**
A: Plugins must be free and open-source. Donations are allowed.

**Q: How often is the registry updated?**
A: The registry syncs every hour and on each release.

**Q: My plugin was rejected, why?**
A: Check the PR comments for review feedback. Common issues: security concerns, missing docs, or quality issues.

---

**Happy Plugin Development! 🚀**

*Built with ❤️ by the StumpfWorks community*
