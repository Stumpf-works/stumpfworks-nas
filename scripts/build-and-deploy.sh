#!/bin/bash
# build-and-deploy.sh - Build and deploy StumpfWorks NAS to APT repository
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APT_SERVER="root@46.4.25.15"
REPO_BASE="/var/www/apt-repo"
GIT_REPO="https://github.com/Stumpf-works/stumpfworks-nas.git"
BRANCH="development"

# Print header
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  StumpfWorks NAS - Build & Deploy System${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Ask for repository target
echo -e "${YELLOW}Select deployment target:${NC}"
echo "  1) Development (testing & development)"
echo "  2) Stable (production releases)"
echo ""
read -p "Enter choice [1-2]: " REPO_CHOICE

case $REPO_CHOICE in
    1)
        REPO_TYPE="development"
        REPO_PATH="${REPO_BASE}/development"
        echo -e "${YELLOW}→ Deploying to DEVELOPMENT repository${NC}"
        ;;
    2)
        REPO_TYPE="stable"
        REPO_PATH="${REPO_BASE}/stable"
        echo -e "${GREEN}→ Deploying to STABLE repository${NC}"

        # Extra confirmation for stable
        echo ""
        echo -e "${RED}⚠️  WARNING: Deploying to STABLE (production)${NC}"
        read -p "Are you sure? [y/N]: " CONFIRM
        if [[ ! $CONFIRM =~ ^[Yy]$ ]]; then
            echo "Aborted."
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}❌ Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  📡 Connecting to APT Server${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Verify SSH connection
echo "🔑 Testing SSH connection..."
if ! ssh -o ConnectTimeout=10 "$APT_SERVER" "echo '✓ Connected'" 2>/dev/null; then
    echo -e "${RED}❌ Failed to connect to APT server${NC}"
    echo "   Check SSH keys and network connectivity"
    exit 1
fi
echo ""

# Execute build on APT server
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  🔨 Building on APT Server${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

ssh "$APT_SERVER" "export REPO_TYPE='$REPO_TYPE' REPO_PATH='$REPO_PATH' BRANCH='$BRANCH' GIT_REPO='$GIT_REPO'; bash" <<'ENDSSH'
set -e

# Colors for remote output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BUILD_DIR="/tmp/stumpfworks-nas-build"
REPO_TYPE="$REPO_TYPE"
REPO_PATH="$REPO_PATH"

echo -e "${YELLOW}📥 Fetching latest code from GitHub...${NC}"

# Clean previous build
if [ -d "$BUILD_DIR" ]; then
    echo "   Cleaning previous build directory..."
    rm -rf "$BUILD_DIR"
fi

# Clone repository
git clone -b "$BRANCH" --depth 1 "$GIT_REPO" "$BUILD_DIR"
cd "$BUILD_DIR"

echo -e "${GREEN}   ✓ Code fetched successfully${NC}"
echo ""

# Get version (without -dirty since we're working from a clean clone)
VERSION=$(git describe --tags --always 2>/dev/null || echo "0.1.0")
VERSION=${VERSION#v}  # Remove leading 'v'

echo -e "${BLUE}📊 Build Information:${NC}"
echo "   Version: $VERSION"
echo "   Branch: $BRANCH"
echo "   Target Repo: $REPO_TYPE"
echo "   Commit: $(git rev-parse --short HEAD)"
echo ""

# Build frontend FIRST (backend embeds frontend dist/)
echo -e "${YELLOW}🎨 Building frontend...${NC}"
cd frontend
echo "   Running npm ci..."
npm ci --silent

echo "   Running build..."
npm run build
cd ..
echo -e "${GREEN}   ✓ Frontend built successfully${NC}"
echo ""

# Copy frontend dist to backend/embedfs/dist
echo -e "${YELLOW}📁 Copying frontend to backend embed directory...${NC}"
rm -rf backend/embedfs/dist
mkdir -p backend/embedfs
cp -r frontend/dist backend/embedfs/
echo -e "${GREEN}   ✓ Frontend copied${NC}"
echo ""

# Build backend (embeds frontend dist/)
echo -e "${YELLOW}🔨 Building backend...${NC}"
cd backend
echo "   Running go mod download..."
go mod download

echo "   Building for amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-amd64 \
    ./cmd/stumpfworks-server

echo "   Building for arm64..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ../dist/stumpfworks-server-arm64 \
    ./cmd/stumpfworks-server

cd ..
echo -e "${GREEN}   ✓ Backend built successfully${NC}"
echo ""

# Create Debian packages
echo -e "${YELLOW}📦 Creating Debian packages...${NC}"

# Build for amd64
echo "   Building package for amd64..."
bash scripts/build-deb.sh "$VERSION" amd64

# Build for arm64
echo "   Building package for arm64..."
bash scripts/build-deb.sh "$VERSION" arm64

echo -e "${GREEN}   ✓ Packages created successfully${NC}"
echo ""

# Deploy to repository
echo -e "${YELLOW}📤 Deploying to $REPO_TYPE repository...${NC}"

# Ensure repository directories exist
mkdir -p "$REPO_PATH/pool/main"
mkdir -p "$REPO_PATH/dists/$REPO_TYPE/main/binary-amd64"
mkdir -p "$REPO_PATH/dists/$REPO_TYPE/main/binary-arm64"

# Copy packages
echo "   Copying packages to repository..."
cp dist/stumpfworks-nas_${VERSION}_amd64.deb "$REPO_PATH/pool/main/"
cp dist/stumpfworks-nas_${VERSION}_arm64.deb "$REPO_PATH/pool/main/"

echo -e "${GREEN}   ✓ Packages deployed${NC}"
echo ""

# Update repository metadata
echo -e "${YELLOW}🔄 Updating repository metadata...${NC}"

cd "$REPO_PATH"

# Generate Packages files
for ARCH in amd64 arm64; do
    echo "   Generating Packages file for $ARCH..."
    dpkg-scanpackages --arch $ARCH pool/ > dists/$REPO_TYPE/main/binary-$ARCH/Packages
    gzip -kf dists/$REPO_TYPE/main/binary-$ARCH/Packages
done

# Generate Release file
echo "   Generating Release file..."
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
    HASH_CMD=$(echo $HASH | tr '[:upper:]' '[:lower:]' | sed 's/sum$/sum/')
    echo "${HASH}:" >> Release
    find main -type f | while read file; do
        HASH_VALUE=$($HASH_CMD "$file" | cut -d' ' -f1)
        SIZE=$(stat -c%s "$file")
        printf " %s %8d %s\n" "$HASH_VALUE" "$SIZE" "$file" >> Release
    done
done

cd "$REPO_PATH"

echo -e "${GREEN}   ✓ Repository metadata updated${NC}"
echo ""

# Cleanup
echo -e "${YELLOW}🧹 Cleaning up build directory...${NC}"
rm -rf "$BUILD_DIR"
echo -e "${GREEN}   ✓ Cleanup complete${NC}"
echo ""

# Summary
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  ✅ Build & Deployment Complete!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}📦 Deployed:${NC}"
echo "   Repository: $REPO_TYPE"
echo "   Version: $VERSION"
echo "   Architectures: amd64, arm64"
echo ""
echo -e "${GREEN}🌐 Available at:${NC}"
echo "   http://apt.stumpfworks.de/dists/$REPO_TYPE/"
echo ""
echo -e "${GREEN}📥 Install on client systems:${NC}"
if [ "$REPO_TYPE" = "development" ]; then
    echo "   # Add development repository"
    echo "   echo 'deb http://apt.stumpfworks.de development main' | sudo tee /etc/apt/sources.list.d/stumpfworks-dev.list"
else
    echo "   # Add stable repository"
    echo "   echo 'deb http://apt.stumpfworks.de stable main' | sudo tee /etc/apt/sources.list.d/stumpfworks.list"
fi
echo "   sudo apt update"
echo "   sudo apt install stumpfworks-nas"
echo ""

ENDSSH

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  🎉 All Done!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Next steps:"
echo "  1. Test installation on 192.168.178.41"
if [ "$REPO_TYPE" = "development" ]; then
    echo "  2. After successful testing, run this script again with 'Stable' to promote to production"
fi
echo ""
