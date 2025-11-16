#!/bin/bash
# Revision: 2025-11-16 | Author: Claude | Version: 1.0.0
# Build script for creating release binaries for all platforms

set -e

VERSION=${VERSION:-$(git describe --tags --always 2>/dev/null || echo "dev")}
OUTPUT_DIR=${OUTPUT_DIR:-"dist/releases"}

echo "🚀 Building StumpfWorks NAS v$VERSION"
echo "======================================"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Platforms to build
declare -A PLATFORMS=(
    ["linux/amd64"]="stumpfworks-nas-linux-amd64"
    ["linux/arm64"]="stumpfworks-nas-linux-arm64"
    ["linux/arm/7"]="stumpfworks-nas-linux-armv7"
    ["darwin/amd64"]="stumpfworks-nas-darwin-amd64"
    ["darwin/arm64"]="stumpfworks-nas-darwin-arm64"
)

# Build for each platform
for platform in "${!PLATFORMS[@]}"; do
    IFS='/' read -r GOOS GOARCH GOARM <<< "$platform"
    output="${PLATFORMS[$platform]}"

    echo "📦 Building for $GOOS/$GOARCH${GOARM:+v$GOARM}..."

    cd backend

    if [ -n "$GOARM" ]; then
        GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build \
            -ldflags="-s -w -X main.AppVersion=$VERSION" \
            -o "../$OUTPUT_DIR/$output" \
            cmd/stumpfworks-server/main.go
    else
        GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags="-s -w -X main.AppVersion=$VERSION" \
            -o "../$OUTPUT_DIR/$output" \
            cmd/stumpfworks-server/main.go
    fi

    cd ..

    size=$(du -h "$OUTPUT_DIR/$output" | cut -f1)
    echo "   ✅ Built $output ($size)"
done

# Build frontend
echo ""
echo "🎨 Building frontend..."
cd frontend
npm run build
cd ..
cp -r frontend/dist "$OUTPUT_DIR/frontend"
echo "   ✅ Frontend built"

# Create checksums
echo ""
echo "🔐 Creating checksums..."
cd "$OUTPUT_DIR"
sha256sum stumpfworks-nas-* > checksums.txt
echo "   ✅ Checksums created"
cd ../..

# Summary
echo ""
echo "✅ Release build complete!"
echo "======================================"
echo "📂 Output directory: $OUTPUT_DIR"
echo ""
ls -lh "$OUTPUT_DIR"/stumpfworks-nas-* 2>/dev/null || true
echo ""
echo "📋 Files:"
find "$OUTPUT_DIR" -type f -name "stumpfworks-nas-*" -o -name "checksums.txt" | sed 's/^/   /'
