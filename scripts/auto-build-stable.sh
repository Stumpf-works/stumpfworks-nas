#!/bin/bash

# Auto-build script for StumpfWorks NAS (Stable/Production Repository)
# This script:
# - Checks for new commits on main branch
# - Builds only if there are changes
# - Deploys to stable repository
# - Sends Discord notifications
# - Runs daily at 20:00 via cron

set -e

# Configuration
REPO_URL="https://github.com/Stumpf-works/stumpfworks-nas.git"
BRANCH="main"
BUILD_DIR="/tmp/stumpfworks-nas-build-stable"
STATE_FILE="/var/lib/stumpfworks-nas/auto-build-stable-state"
LOG_FILE="/var/log/stumpfworks-nas/auto-build-stable.log"
DISCORD_WEBHOOK="https://discord.com/api/webhooks/1444087366410305739/j4dYH00dtK0fAgnAD4U2PgeXZKVcqnvDYzW4I6h1-EoBJXhtOUX2H8yo6nyAvwegf3GG"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Ensure directories exist
mkdir -p "$(dirname "$STATE_FILE")"
mkdir -p "$(dirname "$LOG_FILE")"

# Logging function
log() {
    echo -e "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Discord notification function
send_discord() {
    local title="$1"
    local description="$2"
    local color="$3"  # Decimal color code
    local fields="$4"  # JSON array of fields

    local json_payload=$(cat <<EOF
{
  "embeds": [{
    "title": "$title",
    "description": "$description",
    "color": $color,
    "fields": $fields,
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%S.000Z)",
    "footer": {
      "text": "StumpfWorks NAS Auto-Build (Stable)"
    }
  }]
}
EOF
)

    curl -H "Content-Type: application/json" \
         -d "$json_payload" \
         "$DISCORD_WEBHOOK" \
         -s -o /dev/null
}

# Error handler
handle_error() {
    local exit_code=$?
    log "${RED}❌ Build failed with exit code $exit_code${NC}"

    send_discord \
        "❌ Stable Build Failed" \
        "The automated stable build encountered an error." \
        15158332 \
        "[{\"name\": \"Exit Code\", \"value\": \"$exit_code\", \"inline\": true}, {\"name\": \"Branch\", \"value\": \"$BRANCH\", \"inline\": true}]"

    # Cleanup
    rm -rf "$BUILD_DIR"
    exit $exit_code
}

trap handle_error ERR

# Start
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log "${BLUE}  StumpfWorks NAS - Auto Build (Stable/Production)${NC}"
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Check for new commits
log "${YELLOW}Checking for new commits on $BRANCH...${NC}"

# Get latest commit hash from GitHub
LATEST_COMMIT=$(git ls-remote "$REPO_URL" "refs/heads/$BRANCH" | awk '{print $1}')
log "Latest commit on GitHub: $LATEST_COMMIT"

# Get last built commit
LAST_BUILT_COMMIT=""
if [ -f "$STATE_FILE" ]; then
    LAST_BUILT_COMMIT=$(cat "$STATE_FILE")
    log "Last built commit: $LAST_BUILT_COMMIT"
fi

# Check if we need to build
if [ "$LATEST_COMMIT" = "$LAST_BUILT_COMMIT" ]; then
    log "${GREEN}✓ No new commits - skipping build${NC}"
    exit 0
fi

log "${YELLOW}New commits detected - starting build...${NC}"
log ""

# Send build start notification
COMMIT_SHORT="${LATEST_COMMIT:0:7}"
send_discord \
    "🔨 Stable Build Started" \
    "Building new stable release from main branch." \
    16753920 \
    "[{\"name\": \"Branch\", \"value\": \"$BRANCH\", \"inline\": true}, {\"name\": \"Commit\", \"value\": \"$COMMIT_SHORT\", \"inline\": true}]"

# Cleanup old build directory
rm -rf "$BUILD_DIR"

# Clone repository
log "${YELLOW}📥 Fetching latest code from GitHub...${NC}"
git clone -b "$BRANCH" --depth 1 "$REPO_URL" "$BUILD_DIR"
cd "$BUILD_DIR"

log "${GREEN}✓ Code fetched successfully${NC}"
log ""

# Get version from tags (stable uses direct tag version without -dev suffix)
git fetch --unshallow 2>/dev/null || true
git fetch --tags 2>/dev/null || true

# Get the latest tag (should be on main branch for stable builds)
LATEST_TAG=$(git describe --tags --exact-match 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "")
COMMIT_SHORT=$(git rev-parse --short HEAD)

if [ -n "$LATEST_TAG" ]; then
    # Remove 'v' prefix if present
    VERSION=${LATEST_TAG#v}
else
    # Fallback: use commit hash if no tag exists
    VERSION="0.1.0+${COMMIT_SHORT}"
fi

log "${BLUE}📊 Build Information:${NC}"
log "   Version: $VERSION"
log "   Branch: $BRANCH"
log "   Target Repo: stable"
log "   Commit: $COMMIT_SHORT"
log ""

# Build frontend
log "${YELLOW}🎨 Building frontend...${NC}"
cd frontend
npm ci --silent
npm run build
log "${GREEN}✓ Frontend built successfully${NC}"
log ""

# Copy frontend to backend embed directory
log "${YELLOW}📁 Copying frontend to backend embed directory...${NC}"
cd ..
rm -rf backend/embedfs/dist
mkdir -p backend/embedfs
cp -r frontend/dist backend/embedfs/
log "${GREEN}✓ Frontend copied${NC}"
log ""

# Build backend
log "${YELLOW}🔨 Building backend...${NC}"
cd backend

# Build for multiple architectures
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-amd64 \
    ./cmd/stumpfworks-server

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-arm64 \
    ./cmd/stumpfworks-server

cd ..
log "${GREEN}✓ Backend built successfully${NC}"
log ""

# Create Debian packages
log "${YELLOW}📦 Creating Debian packages...${NC}"
bash scripts/build-deb.sh
log "${GREEN}✓ Packages created successfully${NC}"
log ""

# Deploy to stable repository
log "${YELLOW}📤 Deploying to stable repository...${NC}"

# Copy packages to shared pool
cp dist/*.deb /var/www/apt-repo/pool/main/

# Update stable repository metadata
cd /var/www/apt-repo/dists/stable/main/binary-amd64
dpkg-scanpackages --multiversion . > Packages
gzip -k -f Packages

cd /var/www/apt-repo/dists/stable/main/binary-arm64
dpkg-scanpackages --multiversion . > Packages
gzip -k -f Packages

# Generate Release file for stable
cd /var/www/apt-repo/dists/stable
cat > Release <<EOF
Origin: StumpfWorks
Label: StumpfWorks NAS Stable
Suite: stable
Codename: stable
Architectures: amd64 arm64
Components: main
Description: StumpfWorks NAS Stable Repository
Date: $(date -Ru)
EOF

# Add checksums to Release file
echo "MD5Sum:" >> Release
find main -type f \( -name "Packages" -o -name "Packages.gz" \) -exec md5sum {} \; | sed 's|main/| |' >> Release

echo "SHA1:" >> Release
find main -type f \( -name "Packages" -o -name "Packages.gz" \) -exec sha1sum {} \; | sed 's|main/| |' >> Release

echo "SHA256:" >> Release
find main -type f \( -name "Packages" -o -name "Packages.gz" \) -exec sha256sum {} \; | sed 's|main/| |' >> Release

log "${GREEN}✓ Packages deployed${NC}"
log ""

# Update repository metadata
log "${YELLOW}🔄 Updating repository metadata...${NC}"
# Repository metadata is already updated above
log "${GREEN}✓ Repository metadata updated${NC}"
log ""

# Save state
echo "$LATEST_COMMIT" > "$STATE_FILE"

# Cleanup
log "${YELLOW}🧹 Cleaning up build directory...${NC}"
rm -rf "$BUILD_DIR"
log "${GREEN}✓ Cleanup complete${NC}"
log ""

log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log "${GREEN}✅ Auto-Build Complete!${NC}"
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log ""
log "${GREEN}📦 Deployed:${NC}"
log "   Repository: stable"
log "   Version: $VERSION"
log "   Architectures: amd64, arm64"
log "   Commit: $LATEST_COMMIT"
log ""
log "${GREEN}🌐 Available at:${NC}"
log "   http://apt.stumpfworks.de/dists/stable/"
log ""

# Send success notification
send_discord \
    "✅ Stable Build Successful" \
    "New stable release deployed successfully!" \
    3066993 \
    "[{\"name\": \"Version\", \"value\": \"$VERSION\", \"inline\": true}, {\"name\": \"Commit\", \"value\": \"$COMMIT_SHORT\", \"inline\": true}, {\"name\": \"Repository\", \"value\": \"stable\", \"inline\": true}]"
