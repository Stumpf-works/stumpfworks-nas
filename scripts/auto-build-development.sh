#!/bin/bash
# auto-build-development.sh - Automatically build and deploy development versions

# Configuration
REPO_URL="https://github.com/Stumpf-works/stumpfworks-nas.git"
BRANCH="development"
REPO_BASE="/var/www/apt-repo"
REPO_TYPE="development"
REPO_PATH="${REPO_BASE}/dists/development"
STATE_FILE="/var/lib/stumpfworks-nas/last-build-commit"
LOG_FILE="/var/log/stumpfworks-nas/auto-build.log"
DISCORD_WEBHOOK="https://discord.com/api/webhooks/1444087366410305739/j4dYH00dtK0fAgnAD4U2PgeXZKVcqnvDYzW4I6h1-EoBJXhtOUX2H8yo6nyAvwegf3GG"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Ensure directories exist
mkdir -p "$(dirname $STATE_FILE)"
mkdir -p "$(dirname $LOG_FILE)"

# Log function
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
      "text": "StumpfWorks NAS Auto-Build"
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
    local line_number=$1

    log "${RED}❌ Build failed at line $line_number with exit code $exit_code${NC}"

    # Send Discord notification - Build failed
    send_discord \
        "❌ Build Failed" \
        "Build process failed at line $line_number" \
        15158332 \
        "[{\"name\":\"Exit Code\",\"value\":\"\`$exit_code\`\",\"inline\":true},{\"name\":\"Log\",\"value\":\"Check \`$LOG_FILE\`\",\"inline\":true}]"

    exit $exit_code
}

trap 'handle_error $LINENO' ERR
set -e

log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log "${BLUE}  StumpfWorks NAS - Auto Build (Development)${NC}"
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Get latest commit from GitHub
log "${YELLOW}Checking for new commits...${NC}"
LATEST_COMMIT=$(git ls-remote "$REPO_URL" "refs/heads/$BRANCH" | cut -f1)

if [ -z "$LATEST_COMMIT" ]; then
    log "${RED}❌ Failed to fetch latest commit from GitHub${NC}"
    exit 1
fi

log "Latest commit on GitHub: $LATEST_COMMIT"

# Read last built commit
LAST_COMMIT=""
if [ -f "$STATE_FILE" ]; then
    LAST_COMMIT=$(cat "$STATE_FILE")
    log "Last built commit: $LAST_COMMIT"
else
    log "${YELLOW}No previous build found${NC}"
fi

# Check if we need to build
if [ "$LATEST_COMMIT" = "$LAST_COMMIT" ]; then
    log "${GREEN}✓ No new commits - skipping build${NC}"
    log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    exit 0
fi

log "${YELLOW}New commits detected - starting build...${NC}"
log ""

# Send Discord notification - Build started
COMMIT_SHORT=$(echo $LATEST_COMMIT | cut -c1-7)
send_discord \
    "🔨 Build Started" \
    "New commits detected in development branch" \
    15844367 \
    "[{\"name\":\"Commit\",\"value\":\"\`$COMMIT_SHORT\`\",\"inline\":true},{\"name\":\"Branch\",\"value\":\"\`$BRANCH\`\",\"inline\":true}]"

# Build using the existing script
BUILD_DIR="/tmp/stumpfworks-nas-build"
export REPO_TYPE REPO_PATH BRANCH GIT_REPO="$REPO_URL"

# Clean previous build
if [ -d "$BUILD_DIR" ]; then
    log "Cleaning previous build directory..."
    rm -rf "$BUILD_DIR"
fi

# Clone repository
log "${YELLOW}📥 Fetching latest code from GitHub...${NC}"
git clone -b "$BRANCH" --depth 1 "$REPO_URL" "$BUILD_DIR"
cd "$BUILD_DIR"

log "${GREEN}✓ Code fetched successfully${NC}"
log ""

# Get version - use latest tag from any branch
git fetch --tags 2>/dev/null || true
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
COMMIT_SHORT=$(git rev-parse --short HEAD)

if [ -n "$LATEST_TAG" ]; then
    # Remove 'v' prefix if present
    BASE_VERSION=${LATEST_TAG#v}
    # For development builds, add -dev suffix and commit hash
    VERSION="${BASE_VERSION}-dev+${COMMIT_SHORT}"
else
    # Fallback if no tags exist
    VERSION="0.1.0-dev+${COMMIT_SHORT}"
fi

log "${BLUE}📊 Build Information:${NC}"
log "   Version: $VERSION"
log "   Branch: $BRANCH"
log "   Target Repo: $REPO_TYPE"
log "   Commit: $(git rev-parse --short HEAD)"
log ""

# Build frontend
log "${YELLOW}🎨 Building frontend...${NC}"
cd frontend
npm ci --silent >> "$LOG_FILE" 2>&1
npm run build >> "$LOG_FILE" 2>&1
cd ..
log "${GREEN}✓ Frontend built successfully${NC}"
log ""

# Copy frontend to backend embed directory
log "${YELLOW}📁 Copying frontend to backend embed directory...${NC}"
rm -rf backend/embedfs/dist
mkdir -p backend/embedfs
cp -r frontend/dist backend/embedfs/
log "${GREEN}✓ Frontend copied${NC}"
log ""

# Build backend
log "${YELLOW}🔨 Building backend...${NC}"
cd backend
go mod download >> "$LOG_FILE" 2>&1

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-amd64 \
    ./cmd/stumpfworks-server >> "$LOG_FILE" 2>&1

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-arm64 \
    ./cmd/stumpfworks-server >> "$LOG_FILE" 2>&1

cd ..
log "${GREEN}✓ Backend built successfully${NC}"
log ""

# Create Debian packages
log "${YELLOW}📦 Creating Debian packages...${NC}"
bash scripts/build-deb.sh "$VERSION" amd64 >> "$LOG_FILE" 2>&1
bash scripts/build-deb.sh "$VERSION" arm64 >> "$LOG_FILE" 2>&1
log "${GREEN}✓ Packages created successfully${NC}"
log ""

# Deploy to repository
log "${YELLOW}📤 Deploying to development repository...${NC}"

# Ensure repository directories exist
mkdir -p "$REPO_PATH/main/binary-amd64"
mkdir -p "$REPO_PATH/main/binary-arm64"
mkdir -p "$REPO_BASE/pool/main"

# Copy packages to shared pool
cp dist/stumpfworks-nas_${VERSION}_amd64.deb "$REPO_BASE/pool/main/"
cp dist/stumpfworks-nas_${VERSION}_arm64.deb "$REPO_BASE/pool/main/"

log "${GREEN}✓ Packages deployed${NC}"
log ""

# Update repository metadata
log "${YELLOW}🔄 Updating repository metadata...${NC}"

cd "$REPO_BASE"

# Generate Packages files
for ARCH in amd64 arm64; do
    dpkg-scanpackages --arch $ARCH pool/ > dists/$REPO_TYPE/main/binary-$ARCH/Packages 2>> "$LOG_FILE"
    gzip -kf dists/$REPO_TYPE/main/binary-$ARCH/Packages
done

# Generate Release file
cat > dists/$REPO_TYPE/Release <<EOF
Origin: StumpfWorks
Label: StumpfWorks NAS ${REPO_TYPE^}
Suite: $REPO_TYPE
Codename: $REPO_TYPE
Version: $VERSION
Architectures: amd64 arm64
Components: main
Description: StumpfWorks NAS ${REPO_TYPE^} Repository
Date: $(date -Ru)
EOF

# Generate checksums
cd dists/$REPO_TYPE
for HASH in MD5Sum SHA1 SHA256; do
    HASH_CMD=$(echo $HASH | tr '[:upper:]' '[:lower:]')
    [[ ! $HASH_CMD =~ sum$ ]] && HASH_CMD="${HASH_CMD}sum"

    echo "${HASH}:" >> Release
    find main -type f | while read file; do
        HASH_VALUE=$($HASH_CMD "$file" | cut -d' ' -f1)
        SIZE=$(stat -c%s "$file")
        printf " %s %8d %s\n" "$HASH_VALUE" "$SIZE" "$file" >> Release
    done
done

cd "$REPO_BASE"

log "${GREEN}✓ Repository metadata updated${NC}"
log ""

# Cleanup
log "${YELLOW}🧹 Cleaning up build directory...${NC}"
rm -rf "$BUILD_DIR"
log "${GREEN}✓ Cleanup complete${NC}"
log ""

# Save current commit as last built
echo "$LATEST_COMMIT" > "$STATE_FILE"

# Summary
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log "${GREEN}✅ Auto-Build Complete!${NC}"
log "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
log ""
log "${GREEN}📦 Deployed:${NC}"
log "   Repository: $REPO_TYPE"
log "   Version: $VERSION"
log "   Architectures: amd64, arm64"
log "   Commit: $LATEST_COMMIT"
log ""
log "${GREEN}🌐 Available at:${NC}"
log "   http://apt.stumpfworks.de/dists/$REPO_TYPE/"
log ""

# Send Discord notification - Build successful
COMMIT_SHORT=$(echo $LATEST_COMMIT | cut -c1-7)
send_discord \
    "✅ Build Successful" \
    "New version deployed to development repository" \
    3066993 \
    "[{\"name\":\"Version\",\"value\":\"\`$VERSION\`\",\"inline\":true},{\"name\":\"Commit\",\"value\":\"\`$COMMIT_SHORT\`\",\"inline\":true},{\"name\":\"Repository\",\"value\":\"[Development](http://apt.stumpfworks.de/dists/$REPO_TYPE/)\",\"inline\":false}]"
