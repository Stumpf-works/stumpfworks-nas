#!/bin/bash
set -e

VERSION=${1:-$(git describe --tags --always 2>/dev/null || echo "0.1.0")}
ARCH=${2:-"amd64"}  # Default to amd64, or use "all" for multi-arch deployment
SERVER="root@46.4.25.15"
REPO_PATH="/var/www/apt-repo/pool/main/"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  📦 StumpfWorks NAS Deployment v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Determine which architectures to deploy
if [ "$ARCH" = "all" ]; then
    ARCHITECTURES=("amd64" "arm64" "armhf")
    echo "🌐 Multi-architecture deployment"
else
    ARCHITECTURES=("$ARCH")
    echo "🔧 Single architecture deployment: $ARCH"
fi
echo ""

# Deploy each architecture
for DEPLOY_ARCH in "${ARCHITECTURES[@]}"; do
    DEB_FILE="dist/stumpfworks-nas_${VERSION}_${DEPLOY_ARCH}.deb"

    echo "──────────────────────────────────────────────────────────────"
    echo "  Deploying $DEPLOY_ARCH"
    echo "──────────────────────────────────────────────────────────────"

    # Check if .deb exists
    if [ ! -f "$DEB_FILE" ]; then
        echo "❌ Error: $DEB_FILE not found!"
        echo "   Run 'make build-multiarch' first."
        exit 1
    fi

    echo "📤 Uploading $DEB_FILE to APT server..."
    scp "$DEB_FILE" "$SERVER:$REPO_PATH"
    echo "   ✓ Uploaded $(basename $DEB_FILE) ($(du -h "$DEB_FILE" | cut -f1))"
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🔄 Updating Repository Metadata"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

ssh "$SERVER" "update-apt-repo"
echo "   ✓ Repository metadata updated!"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  ✅ Deployment Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 Deployed packages:"
for DEPLOY_ARCH in "${ARCHITECTURES[@]}"; do
    echo "   ✓ stumpfworks-nas_${VERSION}_${DEPLOY_ARCH}.deb"
done
echo ""
echo "🌐 Package available at:"
echo "   http://apt.stumpfworks.de"
echo ""
echo "📥 Install with:"
echo "   sudo apt update"
echo "   sudo apt install stumpfworks-nas"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Verify packages in repository
echo ""
echo "🔍 Verifying packages in repository..."
if ssh "$SERVER" "apt-cache policy stumpfworks-nas" | grep -q "$VERSION"; then
    echo "   ✓ Packages verified successfully!"

    # Show available architectures
    echo ""
    echo "📋 Available architectures:"
    for DEPLOY_ARCH in "${ARCHITECTURES[@]}"; do
        echo "   • $DEPLOY_ARCH"
    done
else
    echo "   ⚠️  Warning: Could not verify package version"
fi
echo ""
