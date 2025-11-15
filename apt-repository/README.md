# StumpfWorks NAS - APT Repository

This directory contains tools and scripts for managing the StumpfWorks NAS APT repository.

## 📋 Overview

The APT repository enables users to install and update StumpfWorks NAS using standard Debian package management tools:

```bash
sudo apt install stumpfworks-nas
sudo apt update && sudo apt upgrade
```

## 🚀 Quick Start

### For Repository Maintainers

#### 1. Initialize the Repository

```bash
cd apt-repository
./init-repository.sh
```

This will:
- Create the repository directory structure
- Generate package index
- Export GPG public key
- Create index.html for GitHub Pages
- Copy any existing .deb packages

#### 2. Add New Package Version

```bash
# Option 1: Add specific package
./update-repository.sh --add ../stumpfworks-nas_1.0.1-1_amd64.deb

# Option 2: Just update with packages already in pool/
./update-repository.sh --update

# Option 3: Run without arguments to update everything
./update-repository.sh
```

#### 3. Deploy to GitHub Pages

**Option A: Separate Repository (Recommended)**

1. Create new repository: `stumpfworks-apt-repo`

2. Copy repository contents:
   ```bash
   cd apt-repository
   cp -r pool dists stumpfworks.gpg index.html /path/to/stumpfworks-apt-repo/
   ```

3. Push to GitHub:
   ```bash
   cd /path/to/stumpfworks-apt-repo
   git init
   git checkout -b gh-pages
   git add .
   git commit -m "Initial APT repository"
   git remote add origin https://github.com/Stumpf-works/stumpfworks-apt-repo.git
   git push -u origin gh-pages
   ```

4. Enable GitHub Pages:
   - Go to repository Settings → Pages
   - Source: `gh-pages` branch
   - URL will be: `https://stumpf-works.github.io/stumpfworks-apt-repo`

**Option B: Same Repository**

```bash
# Push to gh-pages branch of this repo
git checkout --orphan gh-pages
git rm -rf .
cp -r apt-repository/* .
git add .
git commit -m "APT repository"
git push origin gh-pages
```

### For End Users

#### One-Line Installation

```bash
curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/install.sh | sudo bash
```

#### Manual Installation

```bash
# Add GPG key
curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/stumpfworks.gpg | \
  sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks-archive-keyring.gpg

# Add repository
echo "deb [signed-by=/usr/share/keyrings/stumpfworks-archive-keyring.gpg] https://stumpf-works.github.io/stumpfworks-apt-repo stable main" | \
  sudo tee /etc/apt/sources.list.d/stumpfworks-nas.list

# Install
sudo apt update
sudo apt install stumpfworks-nas
```

## 📁 Directory Structure

```
apt-repository/
├── pool/                           # Package pool
│   └── main/                       # Main component
│       └── *.deb                   # .deb packages
├── dists/                          # Distribution metadata
│   └── stable/                     # Stable distribution
│       ├── Release                 # Release file
│       ├── Release.gpg             # Detached signature
│       ├── InRelease               # Inline signed release
│       └── main/                   # Main component
│           └── binary-amd64/       # Architecture
│               ├── Packages        # Package index
│               └── Packages.gz     # Compressed index
├── stumpfworks.gpg                 # Public GPG key
├── index.html                      # Web interface
├── init-repository.sh              # Initialize repository
├── update-repository.sh            # Update repository
└── install-stumpfworks-nas.sh      # User installation script
```

## 🔐 GPG Key Setup

### Generate GPG Key (First Time)

```bash
# Generate new GPG key
gpg --full-generate-key

# Choose:
# - Type: RSA and RSA
# - Size: 4096 bits
# - Expiration: 2y (2 years)
# - Name: Stumpf.Works Team
# - Email: contact@stumpf.works

# Export private key for GitHub Actions (save securely!)
gpg --export-secret-keys contact@stumpf.works | base64
```

### Configure GitHub Secrets

For automatic signing in GitHub Actions, add these secrets:

1. Go to repository Settings → Secrets and variables → Actions
2. Add secrets:
   - `GPG_PRIVATE_KEY`: Your private GPG key (base64 encoded)
   - `APT_REPO_TOKEN`: Personal Access Token with repo permissions

## 🔄 Automated Publishing

### GitHub Actions Workflow

The repository includes a GitHub Actions workflow (`.github/workflows/publish-apt-repo.yml`) that automatically:

1. Builds .deb package on release
2. Uploads .deb to GitHub release
3. Updates APT repository
4. Signs repository with GPG
5. Pushes to GitHub Pages

### Trigger Automatic Publishing

```bash
# Create and push a new release tag
git tag v1.0.1
git push origin v1.0.1

# Create GitHub release
gh release create v1.0.1 \
  --title "Release v1.0.1" \
  --notes "Release notes here"

# GitHub Actions will automatically:
# - Build .deb package
# - Publish to APT repository
# - Update GitHub Pages
```

## 🛠️ Scripts Reference

### `init-repository.sh`

Initializes a new APT repository.

```bash
./init-repository.sh
```

**What it does:**
- Creates directory structure
- Checks for GPG key
- Exports public key
- Generates package index
- Creates Release file
- Signs repository (if GPG key available)
- Creates index.html

### `update-repository.sh`

Updates the repository with new packages.

```bash
# Update repository
./update-repository.sh

# Add new package
./update-repository.sh --add package.deb
./update-repository.sh -a package.deb

# List packages
./update-repository.sh --list
./update-repository.sh -l

# Help
./update-repository.sh --help
```

**What it does:**
- Adds new packages to pool
- Regenerates package index
- Updates Release file
- Re-signs repository

### `install-stumpfworks-nas.sh`

User-facing installation script.

```bash
# Run as root
sudo ./install-stumpfworks-nas.sh

# Or via curl
curl -fsSL https://YOUR-REPO-URL/install.sh | sudo bash
```

**What it does:**
- Checks system compatibility
- Installs dependencies
- Adds GPG key
- Adds APT repository
- Installs StumpfWorks NAS
- Displays next steps

## 📊 Testing the Repository

### Test Locally

```bash
# Serve repository locally
cd apt-repository
python3 -m http.server 8000

# On another terminal/machine, add the repository
echo "deb [trusted=yes] http://YOUR_IP:8000 stable main" | \
  sudo tee /etc/apt/sources.list.d/stumpfworks-test.list

sudo apt update
sudo apt install stumpfworks-nas
```

### Test on GitHub Pages

```bash
# After deploying to GitHub Pages
curl -I https://stumpf-works.github.io/stumpfworks-apt-repo/

# Test package installation
curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/install.sh | sudo bash
```

## 🔍 Troubleshooting

### Repository Not Found

**Problem:** `E: Failed to fetch ... 404 Not Found`

**Solutions:**
1. Check GitHub Pages is enabled
2. Verify URL is correct
3. Wait a few minutes for GitHub Pages to deploy

### GPG Signature Verification Failed

**Problem:** `GPG error: ... The following signatures couldn't be verified`

**Solutions:**
1. Re-add GPG key:
   ```bash
   curl -fsSL https://YOUR-REPO/stumpfworks.gpg | \
     sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks-archive-keyring.gpg
   ```
2. Check GPG key in repository is up to date
3. Re-sign repository: `./update-repository.sh --update`

### Package Not Found

**Problem:** `E: Unable to locate package stumpfworks-nas`

**Solutions:**
1. Run `sudo apt update`
2. Check package is in `pool/main/`
3. Verify Packages file exists: `dists/stable/main/binary-amd64/Packages`
4. Regenerate index: `./update-repository.sh --update`

## 📚 Additional Resources

- [Debian Repository Format](https://wiki.debian.org/DebianRepository/Format)
- [Debian Repository Setup](https://wiki.debian.org/DebianRepository/Setup)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [GPG Documentation](https://www.gnupg.org/documentation/)

## 🆘 Support

- **Issues:** https://github.com/Stumpf-works/stumpfworks-nas/issues
- **Discussions:** https://github.com/Stumpf-works/stumpfworks-nas/discussions
- **Documentation:** https://github.com/Stumpf-works/stumpfworks-nas/wiki

## 📝 License

Same as StumpfWorks NAS - See LICENSE file in the root directory.

---

**Last Updated:** 2025-11-14
**Maintainer:** Stumpf.Works Team
