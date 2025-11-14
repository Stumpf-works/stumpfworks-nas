# StumpfWorks NAS - APT Repository Setup

**Last Updated:** 2025-11-14
**Status:** Planning Document

---

## 📋 OVERVIEW

This document describes how to set up an APT repository for StumpfWorks NAS, enabling users to install and update via `apt`.

**Benefits:**
- ✅ Users can install with: `sudo apt install stumpfworks-nas`
- ✅ Automatic updates with: `sudo apt update && sudo apt upgrade`
- ✅ Dependency resolution handled by apt
- ✅ Professional distribution method
- ✅ Integration with system package manager

---

## 🎯 OPTION 1: GitHub-Hosted Repository (Free, Recommended)

### Advantages:
- Free hosting
- Integrated with GitHub releases
- HTTPS by default
- CDN-backed (fast downloads)

### Setup Steps:

#### 1. Create Repository Structure

```bash
# Create apt-repo directory
mkdir -p apt-repo/pool/main
mkdir -p apt-repo/dists/stable/main/binary-amd64

# Copy .deb packages
cp ../stumpfworks-nas_1.0.0-1_amd64.deb apt-repo/pool/main/
```

#### 2. Generate GPG Key for Signing

```bash
# Generate key (if you don't have one)
gpg --full-generate-key
# Choose: RSA, 4096 bits, no expiration
# Name: Stumpf.Works Team
# Email: contact@stumpf.works

# Export public key
gpg --armor --export contact@stumpf.works > apt-repo/stumpfworks.gpg

# Get key ID
gpg --list-keys --keyid-format SHORT contact@stumpf.works
```

#### 3. Create Package Index

```bash
cd apt-repo

# Generate Packages file
dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages
gzip -k -f dists/stable/main/binary-amd64/Packages

# Create Release file
cd dists/stable
cat > Release << EOF
Origin: StumpfWorks
Label: StumpfWorks NAS
Suite: stable
Codename: stable
Version: 1.0
Architectures: amd64
Components: main
Description: StumpfWorks NAS Official Repository
Date: $(date -Ru)
EOF

# Add checksums
apt-ftparchive release . >> Release

# Sign Release file
gpg --default-key contact@stumpf.works -abs -o Release.gpg Release
gpg --default-key contact@stumpf.works --clearsign -o InRelease Release
```

#### 4. Upload to GitHub Pages

```bash
# Initialize git in apt-repo
cd apt-repo
git init
git add .
git commit -m "Initial APT repository"

# Create gh-pages branch
git checkout -b gh-pages
git remote add origin https://github.com/Stumpf-works/stumpfworks-apt-repo.git
git push -u origin gh-pages
```

#### 5. Client Configuration

Users add the repository:

```bash
# Add GPG key
curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/stumpfworks.gpg | sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks-archive-keyring.gpg

# Add repository
echo "deb [signed-by=/usr/share/keyrings/stumpfworks-archive-keyring.gpg] https://stumpf-works.github.io/stumpfworks-apt-repo stable main" | sudo tee /etc/apt/sources.list.d/stumpfworks-nas.list

# Update and install
sudo apt update
sudo apt install stumpfworks-nas
```

---

## 🎯 OPTION 2: Self-Hosted Repository

### Requirements:
- Web server (Apache/Nginx)
- HTTPS certificate
- Public domain name

### Nginx Configuration:

```nginx
server {
    listen 443 ssl http2;
    server_name apt.stumpf.works;

    ssl_certificate /etc/letsencrypt/live/apt.stumpf.works/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/apt.stumpf.works/privkey.pem;

    root /var/www/apt-repo;

    location / {
        autoindex on;
    }

    location ~ /(.*)/conf/ {
        deny all;
        return 404;
    }
}
```

---

## 🎯 OPTION 3: Third-Party Services

### packagecloud.io (Free for Open Source)

```bash
# Install packagecloud CLI
gem install package_cloud

# Push package
package_cloud push stumpfworks/stumpfworks-nas/debian/bookworm stumpfworks-nas_1.0.0-1_amd64.deb
```

User installation:
```bash
curl -s https://packagecloud.io/install/repositories/stumpfworks/stumpfworks-nas/script.deb.sh | sudo bash
sudo apt install stumpfworks-nas
```

### Gemfury (Paid, $49/month)

```bash
fury push stumpfworks-nas_1.0.0-1_amd64.deb --as=stumpfworks
```

---

## 🔄 UPDATE WORKFLOW

### When releasing a new version:

#### 1. Build new package

```bash
# Update version in:
# - backend/cmd/stumpfworks-server/main.go (AppVersion)
# - frontend/package.json (version)
# - debian/changelog (add new entry)

# Build
dpkg-buildpackage -b -uc -us
```

#### 2. Add to repository

```bash
# Copy new .deb
cp ../stumpfworks-nas_1.0.1-1_amd64.deb apt-repo/pool/main/

# Regenerate Packages file
cd apt-repo
dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages
gzip -k -f dists/stable/main/binary-amd64/Packages

# Update Release file
cd dists/stable
apt-ftparchive release . > Release
gpg --default-key contact@stumpf.works -abs -o Release.gpg Release
gpg --default-key contact@stumpf.works --clearsign -o InRelease Release

# Commit and push
cd ../..
git add .
git commit -m "Release v1.0.1"
git push
```

#### 3. Users get update

```bash
sudo apt update
# → stumpfworks-nas (1.0.1 available)

sudo apt upgrade
# → Upgrading stumpfworks-nas (1.0.0 → 1.0.1)
```

---

## 🔧 AUTOMATION WITH GITHUB ACTIONS

Create `.github/workflows/publish-apt-repo.yml`:

```yaml
name: Publish to APT Repository

on:
  release:
    types: [published]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout main repo
        uses: actions/checkout@v4

      - name: Build .deb package
        run: |
          sudo apt-get update
          sudo apt-get install -y debhelper golang-go nodejs npm
          dpkg-buildpackage -b -uc -us

      - name: Checkout apt-repo
        uses: actions/checkout@v4
        with:
          repository: Stumpf-works/stumpfworks-apt-repo
          path: apt-repo
          ref: gh-pages
          token: ${{ secrets.GH_PAT }}

      - name: Add package to repository
        run: |
          # Copy .deb
          cp ../stumpfworks-nas_*.deb apt-repo/pool/main/

          # Regenerate index
          cd apt-repo
          dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages
          gzip -k -f dists/stable/main/binary-amd64/Packages

          # Update Release
          cd dists/stable
          apt-ftparchive release . > Release

      - name: Sign repository
        env:
          GPG_PRIVATE_KEY: ${{ secrets.GPG_PRIVATE_KEY }}
        run: |
          echo "$GPG_PRIVATE_KEY" | gpg --import
          cd apt-repo/dists/stable
          gpg --default-key contact@stumpf.works -abs -o Release.gpg Release
          gpg --default-key contact@stumpf.works --clearsign -o InRelease Release

      - name: Commit and push
        run: |
          cd apt-repo
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add .
          git commit -m "Add StumpfWorks NAS ${{ github.ref_name }}"
          git push
```

---

## 🌐 WEB UI INTEGRATION

StumpfWorks NAS already has update checking built-in!

### Current Implementation:

**Backend:** `backend/internal/updates/update_service.go`
```go
// Checks GitHub releases for updates
func (s *Service) CheckForUpdates() (*UpdateInfo, error)
```

**Frontend:** Update notifications in Dashboard

### Extend for APT:

Add APT update check in addition to GitHub check:

```go
// Check both GitHub and APT repository
func (s *Service) CheckForUpdates() (*UpdateInfo, error) {
    // 1. Check GitHub releases (current implementation)
    githubUpdate := s.checkGitHub()

    // 2. Check APT repository
    aptUpdate := s.checkAPT()

    // Return whichever is newer
    return selectNewest(githubUpdate, aptUpdate), nil
}

func (s *Service) checkAPT() (*UpdateInfo, error) {
    // Simulate: apt-cache policy stumpfworks-nas
    cmd := exec.Command("apt-cache", "policy", "stumpfworks-nas")
    output, err := cmd.Output()
    // Parse version...
}
```

### Web UI "Update Now" Button:

```typescript
// Frontend: Update button triggers apt upgrade
const handleUpdate = async () => {
  const response = await api.post('/api/v1/system/update', {
    method: 'apt'  // or 'github' for manual download
  });

  // Backend executes: apt-get update && apt-get install stumpfworks-nas
  // Then: systemctl restart stumpfworks-nas
};
```

---

## 📊 REPOSITORY STATISTICS

Track downloads and usage:

### Option 1: GitHub Analytics
- GitHub Pages provides basic traffic stats
- Can see which files are downloaded

### Option 2: Custom Analytics
```nginx
# Nginx: Log .deb downloads
location ~* \.deb$ {
    access_log /var/log/nginx/apt-downloads.log combined;
}
```

### Option 3: PackageCloud Dashboard
- Built-in analytics
- Download counts per version
- Geographic distribution

---

## 🔐 SECURITY BEST PRACTICES

### GPG Key Management:
- ✅ Use 4096-bit RSA key
- ✅ Set expiration date (1-2 years)
- ✅ Store private key securely (GitHub Secrets)
- ✅ Publish public key on keyserver
- ✅ Document key fingerprint in README

### Repository Signing:
- ✅ Always sign Release file
- ✅ Generate InRelease (signed)
- ✅ Use HTTPS for repository access
- ✅ Verify signatures on client side

### Package Verification:
```bash
# Users can verify package signature
dpkg-sig --verify stumpfworks-nas_1.0.0-1_amd64.deb
```

---

## 📚 INSTALLATION VARIANTS

### Variant 1: From APT Repository (Recommended)
```bash
curl -fsSL https://apt.stumpf.works/key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks.gpg
echo "deb [signed-by=/usr/share/keyrings/stumpfworks.gpg] https://apt.stumpf.works stable main" | sudo tee /etc/apt/sources.list.d/stumpfworks.list
sudo apt update
sudo apt install stumpfworks-nas
```

### Variant 2: Direct .deb Download
```bash
wget https://github.com/Stumpf-works/stumpfworks-nas/releases/download/v1.0.0/stumpfworks-nas_1.0.0-1_amd64.deb
sudo dpkg -i stumpfworks-nas_1.0.0-1_amd64.deb
sudo apt-get install -f
```

### Variant 3: From ISO
- Boot from StumpfWorks NAS OS ISO
- Automatic installation
- Updates via APT after installation

---

## 🎯 ROADMAP

### Phase 1: MVP (Week 1-2)
- [ ] Set up GitHub Pages repository
- [ ] Generate GPG key
- [ ] Create initial package index
- [ ] Write client installation script
- [ ] Test apt install/update workflow

### Phase 2: Automation (Week 3)
- [ ] GitHub Actions workflow for auto-publishing
- [ ] Automated signing
- [ ] Version bumping script
- [ ] Release checklist

### Phase 3: Integration (Week 4)
- [ ] Web UI update notifications
- [ ] One-click updates from web UI
- [ ] Update history tracking
- [ ] Rollback support

### Phase 4: Production (Week 5+)
- [ ] Multiple distribution channels (stable, testing, nightly)
- [ ] Mirror repositories
- [ ] Download statistics
- [ ] Security advisories system

---

## 🔗 RESOURCES

**Documentation:**
- Debian Repository HOWTO: https://wiki.debian.org/DebianRepository/Setup
- APT Repository Format: https://wiki.debian.org/DebianRepository/Format
- Reprepro Guide: https://wiki.debian.org/DebianRepository/SetupWithReprepro

**Tools:**
- `reprepro`: Full-featured APT repository manager
- `aptly`: Modern APT repository tool with snapshotting
- `mini-dinstall`: Lightweight repository solution

**Services:**
- PackageCloud: https://packagecloud.io/
- Gemfury: https://gemfury.com/
- CloudSmith: https://cloudsmith.io/

---

## 💡 RECOMMENDATION

**For StumpfWorks NAS v1.0.0:**

Start with **Option 1 (GitHub Pages)** because:
- ✅ Free
- ✅ Simple to set up
- ✅ Integrates with existing GitHub workflow
- ✅ Professional and reliable
- ✅ Easy to automate with GitHub Actions

Later, consider moving to self-hosted or PackageCloud for:
- Multiple distribution channels (stable/testing/nightly)
- Better analytics
- More control

---

**Last Updated:** 2025-11-14
**Status:** Planning - Ready to implement 🚀
