#!/bin/bash
# Build base Debian system using debootstrap

set -e

BUILD_DIR="$1"
DEBIAN_VERSION="$2"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Building base Debian ${DEBIAN_VERSION} system..."

# Run debootstrap to create minimal Debian system
debootstrap \
    --arch=amd64 \
    --variant=minbase \
    --include=systemd,systemd-sysv,udev,kmod \
    "${DEBIAN_VERSION}" \
    "${BUILD_DIR}/chroot" \
    http://deb.debian.org/debian/

echo "✓ Base system created"

# Install additional required packages
echo "Installing system packages..."

# Mount proc, sys, dev for chroot
mount -t proc none "${BUILD_DIR}/chroot/proc"
mount -t sysfs none "${BUILD_DIR}/chroot/sys"
mount -o bind /dev "${BUILD_DIR}/chroot/dev"
mount -o bind /dev/pts "${BUILD_DIR}/chroot/dev/pts"

# Create package list for chroot installation
cat > "${BUILD_DIR}/chroot/tmp/install-packages.sh" << 'PKGEOF'
#!/bin/bash
export DEBIAN_FRONTEND=noninteractive

# Update package lists
apt-get update

# Install essential packages
apt-get install -y \
    linux-image-amd64 \
    live-boot \
    systemd-sysv \
    network-manager \
    openssh-server \
    sudo \
    curl \
    wget \
    vim \
    nano \
    htop \
    net-tools \
    ethtool \
    iproute2 \
    iputils-ping \
    ca-certificates \
    gnupg \
    lsb-release

# Install bootloader
apt-get install -y \
    grub-pc-bin \
    grub-efi-amd64-bin \
    grub-efi-ia32-bin

# Install live-boot dependencies
apt-get install -y \
    live-boot \
    live-config \
    live-config-systemd

# Clean up
apt-get clean
rm -rf /var/lib/apt/lists/*
PKGEOF

chmod +x "${BUILD_DIR}/chroot/tmp/install-packages.sh"

# Run package installation in chroot
chroot "${BUILD_DIR}/chroot" /tmp/install-packages.sh

# Cleanup
rm "${BUILD_DIR}/chroot/tmp/install-packages.sh"

# Unmount
umount "${BUILD_DIR}/chroot/dev/pts"
umount "${BUILD_DIR}/chroot/dev"
umount "${BUILD_DIR}/chroot/sys"
umount "${BUILD_DIR}/chroot/proc"

echo "✓ Base system packages installed"
