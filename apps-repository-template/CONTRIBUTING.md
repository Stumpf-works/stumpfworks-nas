# Contributing to StumpfWorks NAS Plugins

Thank you for your interest in contributing to the StumpfWorks NAS plugin ecosystem! 🎉

This guide will help you create, test, and submit plugins for inclusion in the official repository.

---

## 📋 Table of Contents

- [Before You Start](#before-you-start)
- [Plugin Requirements](#plugin-requirements)
- [Development Workflow](#development-workflow)
- [Plugin Structure](#plugin-structure)
- [Testing Your Plugin](#testing-your-plugin)
- [Submission Process](#submission-process)
- [Review Criteria](#review-criteria)
- [Best Practices](#best-practices)

---

## 🎯 Before You Start

### Check Existing Plugins

Search existing plugins to avoid duplicates:
- Browse the [plugins directory](./plugins/)
- Check [open pull requests](https://github.com/Stumpf-works/stumpfworks-nas-apps/pulls)
- Search [issues](https://github.com/Stumpf-works/stumpfworks-nas-apps/issues)

### Discuss Your Idea

Open an issue to discuss your plugin idea before starting development:
```
Title: [Plugin Proposal] My Plugin Name
Description: Brief description of what your plugin does
```

This helps ensure your contribution will be accepted and avoids wasted effort.

---

## ✅ Plugin Requirements

Every plugin must have:

### Required Files

- ✅ **plugin.json** - Plugin manifest (required)
- ✅ **README.md** - User documentation
- ✅ **CHANGELOG.md** - Version history
- ✅ **LICENSE** - License file (MIT, Apache 2.0, GPL, etc.)

### Required Metadata

In `plugin.json`:
```json
{
  "id": "com.company.plugin-name",       // Required: Unique ID (reverse domain)
  "name": "Plugin Display Name",         // Required
  "version": "1.0.0",                    // Required: Semantic versioning
  "author": "Your Name",                 // Required
  "description": "Short description",    // Required
  "icon": "🔌",                          // Required: Emoji or icon path
  "category": "utilities",               // Required: See categories below
  "min_nas_version": "0.1.0"            // Required
}
```

### Optional but Recommended

- 📸 **screenshots/** - Visual previews
- 📝 **docs/** - Detailed documentation
- ✅ **tests/** - Automated tests
- 🐳 **docker-compose.yml** - For Docker-based plugins

---

## 🛠️ Development Workflow

### 1. Fork & Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR-USERNAME/stumpfworks-nas-apps.git
cd stumpfworks-nas-apps
```

### 2. Choose a Template

```bash
# Basic plugin (standalone Go/Python/etc.)
cp -r templates/basic-plugin plugins/my-plugin

# Docker-based plugin
cp -r templates/docker-plugin plugins/my-plugin

# Full-stack plugin (Backend + Frontend)
cp -r templates/full-stack-plugin plugins/my-plugin
```

### 3. Customize plugin.json

```bash
cd plugins/my-plugin
nano plugin.json
```

**Important**: Choose a unique plugin ID using reverse domain notation:
```
Good: com.github-username.plugin-name
Bad:  my-plugin
```

### 4. Develop Your Plugin

See the main [Plugin Development Guide](https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/DEVELOPMENT.md) for detailed instructions.

### 5. Test Locally

```bash
# Copy to local StumpfWorks NAS
sudo cp -r plugins/my-plugin /var/lib/stumpfworks/plugins/

# Or use Docker for testing
cd plugins/my-plugin
docker-compose up -d
```

### 6. Create Release Archive

```bash
# From your plugin directory
cd plugins/my-plugin

# Create tar.gz (exclude .git, node_modules, etc.)
tar czf ../../releases/my-plugin-v1.0.0.tar.gz \
  --exclude=".git" \
  --exclude="node_modules" \
  --exclude="*.log" \
  --exclude="tmp" \
  .
```

---

## 📁 Plugin Structure

### Minimal Structure

```
plugins/my-plugin/
├── plugin.json          # Plugin manifest
├── README.md            # Documentation
├── CHANGELOG.md         # Version history
├── LICENSE              # License file
└── my-plugin            # Executable (for native plugins)
```

### Docker Plugin Structure

```
plugins/my-plugin/
├── plugin.json
├── README.md
├── CHANGELOG.md
├── LICENSE
├── docker-compose.yml
├── Dockerfile           # Optional: custom image
└── config/             # Configuration files
```

### Full-Stack Plugin Structure

```
plugins/my-plugin/
├── plugin.json
├── README.md
├── CHANGELOG.md
├── LICENSE
├── docker-compose.yml
├── backend/
│   ├── main.go
│   ├── api/
│   └── Dockerfile
├── frontend/
│   ├── src/
│   ├── package.json
│   └── Dockerfile
└── config/
```

---

## 🧪 Testing Your Plugin

### Manual Testing

1. **Install on real StumpfWorks NAS**
   ```bash
   # Copy plugin
   sudo cp -r plugins/my-plugin /var/lib/stumpfworks/plugins/

   # Test via UI or API
   curl -X POST http://localhost:8080/api/v1/plugins/my-plugin/enable
   ```

2. **Test all features**
   - Installation
   - Configuration
   - Start/Stop
   - Core functionality
   - Uninstallation

3. **Check logs**
   ```bash
   journalctl -u stumpfworks-nas -f | grep my-plugin
   ```

### Automated Testing

If you include tests:
```bash
# Go plugins
cd backend && go test ./...

# JavaScript/TypeScript
cd frontend && npm test

# Docker plugins
docker-compose up -d
docker-compose logs -f
```

---

## 📝 Submission Process

### 1. Prepare Your Submission

Checklist:
- [ ] plugin.json is valid
- [ ] README.md is complete
- [ ] CHANGELOG.md has v1.0.0 entry
- [ ] LICENSE file is included
- [ ] Release archive (.tar.gz) is created
- [ ] Tested on StumpfWorks NAS
- [ ] No security vulnerabilities
- [ ] No hardcoded secrets

### 2. Create Pull Request

```bash
# Commit your changes
git add plugins/my-plugin/
git commit -m "Add: My Plugin - Brief description"

# Push to your fork
git push origin main

# Create PR on GitHub
```

**PR Title Format:**
```
[Plugin] Plugin Name - Brief description
```

**PR Description Template:**
```markdown
## Plugin Information

- **Name**: My Plugin
- **ID**: com.github-username.my-plugin
- **Version**: 1.0.0
- **Category**: utilities

## Description

Brief description of what your plugin does.

## Features

- Feature 1
- Feature 2
- Feature 3

## Testing

Describe how you tested the plugin:
- [ ] Installed successfully
- [ ] Starts without errors
- [ ] Core features work
- [ ] Uninstalls cleanly

## Screenshots

[Optional: Add screenshots]

## Checklist

- [ ] plugin.json is valid
- [ ] README.md included
- [ ] CHANGELOG.md included
- [ ] LICENSE included
- [ ] Release archive created
- [ ] Tested on StumpfWorks NAS
- [ ] No security issues
```

### 3. Automated Checks

Your PR will automatically run:
- ✅ plugin.json validation
- ✅ Required files check
- ✅ Docker syntax validation (if applicable)
- ✅ Security scan
- ✅ Basic functionality test

Fix any issues reported by the automated checks.

### 4. Code Review

A maintainer will review your plugin within 3-5 business days. They may:
- Request changes
- Ask questions
- Suggest improvements

Please be responsive to review feedback.

### 5. Approval & Merge

Once approved:
1. Maintainer merges your PR
2. GitHub Actions creates a release
3. Registry is updated automatically
4. Your plugin goes live! 🎉

---

## 🔍 Review Criteria

Plugins are reviewed based on:

### Functionality
- Does it work as described?
- Are all features implemented?
- Is error handling robust?

### Security
- No known vulnerabilities
- No hardcoded secrets
- Proper input validation
- Safe file operations
- Limited permissions

### Quality
- Clean, readable code
- Proper documentation
- Follows best practices
- No unnecessary dependencies

### User Experience
- Easy to configure
- Clear error messages
- Good performance
- Doesn't break other plugins

### Maintenance
- Active author
- Responds to issues
- Plans for updates
- Open to contributions

---

## 🌟 Best Practices

### 1. Security

```bash
# ❌ BAD: Hardcoded credentials
const password = "my-secret-password"

# ✅ GOOD: Environment variables
const password = process.env.PLUGIN_PASSWORD

# ❌ BAD: Running as root
USER root

# ✅ GOOD: Non-root user
USER plugin
```

### 2. Configuration

```bash
# ❌ BAD: Direct config file editing
echo "setting=value" >> /etc/app.conf

# ✅ GOOD: Use plugin.json config section
{
  "config": {
    "setting": "value"
  }
}
```

### 3. Logging

```go
// ❌ BAD: Print to stdout
fmt.Println("Something happened")

// ✅ GOOD: Structured logging
log.Info().Str("event", "something").Msg("Event occurred")
```

### 4. Dependencies

```dockerfile
# ❌ BAD: Latest tag
FROM node:latest

# ✅ GOOD: Specific version
FROM node:18-alpine

# ❌ BAD: Many dependencies
RUN apt-get install -y git curl wget vim emacs ...

# ✅ GOOD: Only what's needed
RUN apk add --no-cache curl
```

### 5. Documentation

```markdown
# ❌ BAD: Minimal README
# My Plugin

Does stuff.

# ✅ GOOD: Comprehensive README
# My Plugin

## Description
Detailed description...

## Features
- Feature 1
- Feature 2

## Installation
Step-by-step instructions...

## Configuration
Available options...

## Troubleshooting
Common issues and solutions...
```

---

## 📚 Resources

- **Main Documentation**: https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/
- **Plugin Examples**: [./plugins/](./plugins/)
- **Templates**: [./templates/](./templates/)
- **API Reference**: https://docs.stumpfworks.com/api

---

## 💬 Getting Help

- **GitHub Discussions**: https://github.com/Stumpf-works/stumpfworks-nas-apps/discussions
- **Discord**: https://discord.gg/stumpfworks
- **Email**: plugins@stumpfworks.com

---

## 📊 Plugin Categories

Choose one of the following for your `plugin.json`:

- **storage** - Object storage, backup, sync
- **media** - Media servers, streaming
- **communication** - VoIP, chat, email
- **development** - Git, CI/CD, databases
- **monitoring** - Metrics, logging, alerting
- **networking** - DNS, VPN, proxy, firewall
- **productivity** - Tasks, notes, documents
- **security** - Password managers, 2FA
- **utilities** - General tools

---

## 🎉 Thank You!

Your contributions help make StumpfWorks NAS better for everyone!

Questions? Feel free to ask in [Discussions](https://github.com/Stumpf-works/stumpfworks-nas-apps/discussions).

Happy coding! 🚀
