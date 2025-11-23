#!/bin/bash
# build-multiarch.sh - Build StumpfWorks NAS for multiple architectures
set -e

VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")}
# Remove leading 'v' from version if present (Debian requirement)
VERSION=${VERSION#v}
# Use absolute path for BUILD_DIR to avoid issues with cd commands
BUILD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/dist"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  StumpfWorks NAS - Multi-Architecture Builder"
echo "  Version: $VERSION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Architectures to build for
ARCHITECTURES=("amd64" "arm64" "armhf")
GOARCH_MAP=("amd64" "arm64" "arm")

# Check if we're on Linux (needed for fakeroot/dpkg-deb)
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo "⚠️  Warning: Not running on Linux. Package building may fail."
    echo "   Building binaries only..."
    BUILD_PACKAGES=false
else
    BUILD_PACKAGES=true
fi

# Step 1: Build frontend once (shared across all architectures)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Step 1/4: Building Frontend"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend && npm install && cd ..
fi

echo "🔨 Building frontend..."
cd frontend && npm run build && cd ..

echo "📁 Copying frontend files for embedding..."
mkdir -p backend/embedfs
rm -rf backend/embedfs/dist
cp -r frontend/dist backend/embedfs/

echo "✅ Frontend build complete!"
echo ""

# Step 2: Build Go dependencies once
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Step 2/4: Preparing Go Dependencies"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "📦 Downloading Go dependencies..."
cd backend && go mod tidy && go mod download && cd ..

echo "✅ Dependencies ready!"
echo ""

# Step 3: Build binaries for all architectures
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Step 3/4: Building Binaries for All Architectures"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

mkdir -p "$BUILD_DIR"

for i in "${!ARCHITECTURES[@]}"; do
    ARCH="${ARCHITECTURES[$i]}"
    GOARCH="${GOARCH_MAP[$i]}"
    GOARM=""

    # Set GOARM for armhf
    if [ "$GOARCH" = "arm" ]; then
        GOARM="7"
    fi

    echo "──────────────────────────────────────────────────────────────"
    echo "  Building for $ARCH (GOARCH=$GOARCH${GOARM:+, GOARM=$GOARM})"
    echo "──────────────────────────────────────────────────────────────"

    # Build main server binary
    echo "  🔨 Building stumpfworks-server..."
    cd backend
    if [ -n "$GOARM" ]; then
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH GOARM=$GOARM go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfworks-server-${GOARCH}" \
            cmd/stumpfworks-server/main.go
    else
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfworks-server-${GOARCH}" \
            cmd/stumpfworks-server/main.go
    fi
    cd ..
    echo "     ✓ stumpfworks-server-${GOARCH} ($(du -h "$BUILD_DIR/stumpfworks-server-${GOARCH}" | cut -f1))"

    # Build stumpfctl
    echo "  🔨 Building stumpfctl..."
    cd backend
    if [ -n "$GOARM" ]; then
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH GOARM=$GOARM go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfctl-${GOARCH}" \
            cmd/stumpfctl/main.go
    else
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfctl-${GOARCH}" \
            cmd/stumpfctl/main.go
    fi
    cd ..
    echo "     ✓ stumpfctl-${GOARCH} ($(du -h "$BUILD_DIR/stumpfctl-${GOARCH}" | cut -f1))"

    # Build stumpfworks-dbsetup
    echo "  🔨 Building stumpfworks-dbsetup..."
    cd backend
    if [ -n "$GOARM" ]; then
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH GOARM=$GOARM go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" \
            cmd/stumpfworks-dbsetup/main.go
    else
        CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build \
            -ldflags="-s -w" \
            -o "../$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" \
            cmd/stumpfworks-dbsetup/main.go
    fi
    cd ..
    echo "     ✓ stumpfworks-dbsetup-${GOARCH} ($(du -h "$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" | cut -f1))"

    echo ""
done

echo "✅ All binaries built successfully!"
echo ""

# Step 4: Build Debian packages for all architectures
if [ "$BUILD_PACKAGES" = true ]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  Step 4/4: Building Debian Packages"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    for ARCH in "${ARCHITECTURES[@]}"; do
        echo "──────────────────────────────────────────────────────────────"
        echo "  Packaging for $ARCH"
        echo "──────────────────────────────────────────────────────────────"
        ./scripts/build-deb.sh "$VERSION" "$ARCH"
        echo ""
    done

    echo "✅ All packages built successfully!"
else
    echo "⏭️  Skipping package building (not on Linux)"
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  ✅ Multi-Architecture Build Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 Built artifacts:"
echo ""

for i in "${!ARCHITECTURES[@]}"; do
    ARCH="${ARCHITECTURES[$i]}"
    GOARCH="${GOARCH_MAP[$i]}"

    echo "  $ARCH:"
    if [ -f "$BUILD_DIR/stumpfworks-server-${GOARCH}" ]; then
        echo "    ✓ stumpfworks-server-${GOARCH} ($(du -h "$BUILD_DIR/stumpfworks-server-${GOARCH}" | cut -f1))"
    fi
    if [ -f "$BUILD_DIR/stumpfctl-${GOARCH}" ]; then
        echo "    ✓ stumpfctl-${GOARCH} ($(du -h "$BUILD_DIR/stumpfctl-${GOARCH}" | cut -f1))"
    fi
    if [ -f "$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" ]; then
        echo "    ✓ stumpfworks-dbsetup-${GOARCH} ($(du -h "$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" | cut -f1))"
    fi
    if [ -f "$BUILD_DIR/stumpfworks-nas_${VERSION}_${ARCH}.deb" ]; then
        echo "    ✓ stumpfworks-nas_${VERSION}_${ARCH}.deb ($(du -h "$BUILD_DIR/stumpfworks-nas_${VERSION}_${ARCH}.deb" | cut -f1))"
    fi
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📥 To deploy all packages:"
echo "   make deploy-multiarch VERSION=$VERSION"
echo ""
echo "📥 To deploy single architecture:"
echo "   make deploy VERSION=$VERSION ARCH=amd64"
echo "   make deploy VERSION=$VERSION ARCH=arm64"
echo "   make deploy VERSION=$VERSION ARCH=armhf"
echo ""
