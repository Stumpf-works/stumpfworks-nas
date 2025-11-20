#!/bin/bash
set -e

VERSION=${1:-$(git describe --tags --always 2>/dev/null || echo "0.1.0")}
DEB_FILE="dist/stumpfworks-nas_${VERSION}_amd64.deb"
SERVER="root@46.4.25.15"
REPO_PATH="/var/www/apt-repo/pool/main/"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  📦 StumpfWorks NAS Deployment v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if .deb exists
if [ ! -f "$DEB_FILE" ]; then
    echo "❌ Error: $DEB_FILE not found!"
    echo "   Run 'make deb' first."
    exit 1
fi

echo "📤 Uploading $DEB_FILE to APT server..."
scp "$DEB_FILE" "$SERVER:$REPO_PATH"

echo "🔄 Updating repository metadata..."
ssh "$SERVER" "update-apt-repo"

echo "✅ Deployment complete!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Package available at:"
echo "   http://apt.stumpfworks.de"
echo ""
echo "📥 Install with:"
echo "   sudo apt update"
echo "   sudo apt install stumpfworks-nas"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Verify package in repository
echo ""
echo "🔍 Verifying package in repository..."
if ssh "$SERVER" "apt-cache policy stumpfworks-nas" | grep -q "$VERSION"; then
    echo "   ✓ Package verified successfully!"
else
    echo "   ⚠️  Warning: Could not verify package version"
fi
