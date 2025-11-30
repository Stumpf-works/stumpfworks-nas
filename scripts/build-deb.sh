#!/bin/bash
# build-deb.sh - Build Debian package for StumpfWorks NAS
set -e

VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")}
# Remove leading 'v' from version if present (Debian requirement)
VERSION=${VERSION#v}
ARCH=${2:-"amd64"}  # Default to amd64 if not specified
BUILD_DIR="$(pwd)/dist"
DEB_DIR="$BUILD_DIR/debian-build-${ARCH}"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  StumpfWorks NAS - Debian Package Builder"
echo "  Version: $VERSION | Architecture: $ARCH"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check for required tools
echo "🔍 Checking build dependencies..."
for cmd in dpkg-deb fakeroot; do
    if ! command -v $cmd &> /dev/null; then
        echo "❌ Error: $cmd is not installed"
        echo "   Install with: sudo apt install $cmd"
        exit 1
    fi
done
echo "   ✓ All dependencies available"
echo ""

# Clean previous build
if [ -d "$DEB_DIR" ]; then
    echo "🧹 Cleaning previous build..."
    rm -rf "$DEB_DIR"
fi

# Create directory structure
echo "📁 Creating package structure..."
mkdir -p "$DEB_DIR/DEBIAN"
mkdir -p "$DEB_DIR/usr/bin"
mkdir -p "$DEB_DIR/etc/stumpfworks-nas"
mkdir -p "$DEB_DIR/etc/systemd/system"
mkdir -p "$DEB_DIR/usr/share/doc/stumpfworks-nas"
mkdir -p "$DEB_DIR/var/lib/stumpfworks-nas"
mkdir -p "$DEB_DIR/var/log/stumpfworks-nas"

# Determine binary suffix based on architecture
case "$ARCH" in
    amd64)
        GOARCH="amd64"
        ;;
    arm64)
        GOARCH="arm64"
        ;;
    armhf)
        GOARCH="arm"
        ;;
    *)
        echo "❌ Error: Unsupported architecture: $ARCH"
        echo "   Supported: amd64, arm64, armhf"
        exit 1
        ;;
esac

# Copy binaries
echo "📦 Copying binaries (GOARCH=$GOARCH)..."
if [ ! -f "$BUILD_DIR/stumpfworks-server-${GOARCH}" ]; then
    echo "❌ Error: stumpfworks-server-${GOARCH} binary not found"
    echo "   Run 'make build-multiarch' first"
    exit 1
fi

cp "$BUILD_DIR/stumpfworks-server-${GOARCH}" "$DEB_DIR/usr/bin/stumpfworks-nas"
chmod 755 "$DEB_DIR/usr/bin/stumpfworks-nas"

if [ -f "$BUILD_DIR/stumpfctl-${GOARCH}" ]; then
    cp "$BUILD_DIR/stumpfctl-${GOARCH}" "$DEB_DIR/usr/bin/stumpfctl"
    chmod 755 "$DEB_DIR/usr/bin/stumpfctl"
fi

if [ -f "$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" ]; then
    cp "$BUILD_DIR/stumpfworks-dbsetup-${GOARCH}" "$DEB_DIR/usr/bin/stumpfworks-dbsetup"
    chmod 755 "$DEB_DIR/usr/bin/stumpfworks-dbsetup"
fi

# Copy configuration
echo "📝 Copying configuration files..."
cp config.yaml.example "$DEB_DIR/etc/stumpfworks-nas/"
chmod 644 "$DEB_DIR/etc/stumpfworks-nas/config.yaml.example"

# Copy systemd services
cp debian/stumpfworks-nas.service "$DEB_DIR/etc/systemd/system/"
chmod 644 "$DEB_DIR/etc/systemd/system/stumpfworks-nas.service"

cp scripts/bridge-firewall.service "$DEB_DIR/etc/systemd/system/"
chmod 644 "$DEB_DIR/etc/systemd/system/bridge-firewall.service"

# Copy bridge firewall setup script
cp scripts/setup-bridge-firewall.sh "$DEB_DIR/usr/local/bin/"
chmod 755 "$DEB_DIR/usr/local/bin/setup-bridge-firewall.sh"

# Copy documentation
echo "📚 Copying documentation..."
cp README.md "$DEB_DIR/usr/share/doc/stumpfworks-nas/"
cp LICENSE "$DEB_DIR/usr/share/doc/stumpfworks-nas/"
if [ -f CHANGELOG.md ]; then
    cp CHANGELOG.md "$DEB_DIR/usr/share/doc/stumpfworks-nas/"
fi

# Create DEBIAN/control
echo "✍️  Creating control file..."
cat > "$DEB_DIR/DEBIAN/control" <<EOF
Package: stumpfworks-nas
Version: $VERSION
Section: admin
Priority: optional
Architecture: $ARCH
Maintainer: Stumpf.Works Team <contact@stumpf.works>
Depends: postgresql (>= 13), postgresql-client, samba (>= 4.0), smbclient, smartmontools, systemd, sudo
Recommends: nfs-kernel-server, lvm2, mdadm, docker.io | docker-ce, btrfs-progs, zfsutils-linux
Homepage: https://github.com/Stumpf-works/stumpfworks-nas
Description: Modern NAS Management System
 Stumpf.Works NAS is a production-ready Network Attached Storage management
 system designed as an alternative to TrueNAS, Unraid, and Synology DSM.
 .
 Features include user management, SMB/NFS shares, Docker integration,
 SMART monitoring, backup management, 2FA, and a modern web UI.
EOF

# Copy maintainer scripts
echo "📋 Copying maintainer scripts..."
cp debian/postinst "$DEB_DIR/DEBIAN/"
cp debian/prerm "$DEB_DIR/DEBIAN/"
cp debian/postrm "$DEB_DIR/DEBIAN/"
chmod 755 "$DEB_DIR/DEBIAN/postinst"
chmod 755 "$DEB_DIR/DEBIAN/prerm"
chmod 755 "$DEB_DIR/DEBIAN/postrm"

# Build the package
echo ""
echo "🔨 Building Debian package..."
DEB_FILE="$BUILD_DIR/stumpfworks-nas_${VERSION}_${ARCH}.deb"
fakeroot dpkg-deb --build "$DEB_DIR" "$DEB_FILE"

# Verify the package
echo ""
echo "🔍 Verifying package..."
dpkg-deb -I "$DEB_FILE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  ✅ Debian package built successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 Package: $DEB_FILE"
echo "📊 Size: $(du -h "$DEB_FILE" | cut -f1)"
echo ""
echo "To install locally:"
echo "  sudo apt install $DEB_FILE"
echo ""
echo "To deploy to repository:"
echo "  make deploy VERSION=$VERSION"
echo ""
