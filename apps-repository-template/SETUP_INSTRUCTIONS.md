# Setup Instructions for stumpfworks-nas-apps Repository

## 🎯 Goal

Create the official plugin repository at: https://github.com/Stumpf-works/stumpfworks-nas-apps

## 📦 Repository Contents

All files for the apps repository are in: `/home/user/stumpfworks-nas-apps-repo/`

### Structure

```
stumpfworks-nas-apps-repo/
├── README.md                                    # Main repository README
├── CONTRIBUTING.md                              # Contribution guidelines
├── CLAUDE_CODE_MASTER_PROMPT.md                # Master prompt for Claude Code
├── SETUP_INSTRUCTIONS.md                        # This file
├── registry.json                                # Plugin registry (auto-generated)
│
├── .github/
│   └── workflows/
│       ├── validate-plugins.yml                 # CI for validation
│       └── update-registry.yml                  # Auto-update registry
│
├── scripts/
│   ├── generate-registry.py                     # Generate registry.json
│   └── validate-plugins.py                      # Validate plugins
│
└── templates/
    └── docker-plugin/                           # Plugin template
        ├── plugin.json
        ├── README_TEMPLATE.md
        └── CHANGELOG.md
```

## 🚀 Setup Steps

### 1. Create GitHub Repository

```bash
# On GitHub:
# 1. Go to https://github.com/Stumpf-works
# 2. Click "New Repository"
# 3. Name: stumpfworks-nas-apps
# 4. Description: Official plugin repository for StumpfWorks NAS
# 5. Public repository
# 6. Click "Create repository"
```

### 2. Initialize Local Repository

```bash
cd /home/user/stumpfworks-nas-apps-repo

# Initialize git
git init

# Add all files
git add .

# First commit
git commit -m "Initial commit: StumpfWorks NAS Apps repository structure"

# Add remote
git remote add origin https://github.com/Stumpf-works/stumpfworks-nas-apps.git

# Push to main branch
git branch -M main
git push -u origin main
```

### 3. Add Asterisk Plugin

```bash
# Copy the asterisk-voip plugin from main repo
mkdir -p plugins/asterisk-voip
cp -r /home/user/stumpfworks-nas/plugins/asterisk-voip/* plugins/asterisk-voip/

# Create release archive
cd plugins/asterisk-voip
tar czf ../../releases/asterisk-voip-v1.0.0-beta.tar.gz \
  --exclude=".git" \
  --exclude="node_modules" \
  --exclude="*.log" \
  .

cd ../..

# Commit
git add plugins/asterisk-voip releases/
git commit -m "Add: Asterisk VoIP Plugin v1.0.0-beta"
git push
```

### 4. Create GitHub Release

```bash
# Tag the release
git tag -a asterisk-voip-v1.0.0-beta -m "Asterisk VoIP Plugin v1.0.0-beta"
git push origin asterisk-voip-v1.0.0-beta

# On GitHub:
# 1. Go to Releases
# 2. Click "Draft a new release"
# 3. Choose tag: asterisk-voip-v1.0.0-beta
# 4. Title: Asterisk VoIP Plugin v1.0.0-beta
# 5. Upload: releases/asterisk-voip-v1.0.0-beta.tar.gz
# 6. Click "Publish release"
```

### 5. Update registry.json

```bash
# Generate registry
python3 scripts/generate-registry.py

# Commit updated registry
git add registry.json
git commit -m "chore: update registry.json"
git push
```

### 6. Enable GitHub Actions

GitHub Actions workflows are already in `.github/workflows/`:
- **validate-plugins.yml** - Runs on every PR/push to validate plugins
- **update-registry.yml** - Auto-updates registry.json daily or on push

No additional setup needed - they will run automatically!

### 7. Test Registry from StumpfWorks NAS

```bash
# On your StumpfWorks NAS instance:

# List available plugins
curl http://localhost:8080/api/v1/store/plugins

# Sync registry
curl -X POST http://localhost:8080/api/v1/store/sync \
  -H "Authorization: Bearer $TOKEN"

# Install Asterisk plugin
curl -X POST http://localhost:8080/api/v1/store/plugins/com.stumpfworks.asterisk-voip/install \
  -H "Authorization: Bearer $TOKEN"
```

## 📝 Adding New Plugins

### For Repository Maintainers

1. Review Pull Request
2. Merge if approved
3. Tag release: `git tag plugin-name-v1.0.0 && git push --tags`
4. GitHub Actions automatically updates registry.json

### For Plugin Developers

1. Fork repository
2. Copy template: `cp -r templates/docker-plugin plugins/my-plugin`
3. Develop plugin
4. Create PR
5. Pass CI checks
6. Await review

See [CONTRIBUTING.md](./CONTRIBUTING.md) for full details.

## 🎨 Using Claude Code with Master Prompt

When starting a Claude Code session to develop a plugin:

1. Reference the master prompt:
   ```
   @CLAUDE_CODE_MASTER_PROMPT.md I want to create a new plugin for...
   ```

2. Claude will have full context about:
   - StumpfWorks NAS architecture
   - Plugin development patterns
   - Best practices
   - Required file structure
   - Testing procedures

## 🤖 Automated Workflows

### Validation (On PR/Push)

- Validates plugin.json syntax
- Checks required files
- Validates docker-compose.yml
- Security scan (TODO)

### Registry Update (Daily + On Push)

- Scans all plugins/*/plugin.json
- Generates registry.json
- Commits if changed

## 📚 Documentation Links

- **Main NAS Repository**: https://github.com/Stumpf-works/stumpfworks-nas
- **Plugin Development Guide**: https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/DEVELOPMENT.md
- **Architecture Overview**: [APPS_REPOSITORY_STRUCTURE.md](https://github.com/Stumpf-works/stumpfworks-nas/blob/main/plugins/APPS_REPOSITORY_STRUCTURE.md)

## ✅ Checklist

- [ ] GitHub repository created
- [ ] Initial commit pushed
- [ ] Asterisk plugin added
- [ ] First release created
- [ ] registry.json generated
- [ ] GitHub Actions enabled
- [ ] Tested from StumpfWorks NAS

## 🎉 Done!

Your plugin repository is now live! Users can:
- Browse plugins in the Store
- Install with one click
- Receive automatic updates
- Discover new plugins as they're added

## 🚀 Next Steps

1. **Add more plugins** - Start with popular services
2. **Promote** - Announce in community forums
3. **Documentation** - Expand wiki with tutorials
4. **Community** - Encourage contributions

---

**Questions?** Open an issue or discussion in the repository!
