#!/bin/bash
# Install StumpfWorks NAS into the chroot

set -e

BUILD_DIR="$1"
VERSION="$2"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "Installing StumpfWorks NAS v${VERSION}..."

# Mount filesystems for chroot
mount -t proc none "${BUILD_DIR}/chroot/proc"
mount -t sysfs none "${BUILD_DIR}/chroot/sys"
mount -o bind /dev "${BUILD_DIR}/chroot/dev"
mount -o bind /dev/pts "${BUILD_DIR}/chroot/dev/pts"

# Option 1: Install from pre-built .deb package (if exists)
if [ -f "${PROJECT_ROOT}/../stumpfworks-nas_${VERSION}-1_amd64.deb" ]; then
    echo "Installing from .deb package..."
    cp "${PROJECT_ROOT}/../stumpfworks-nas_${VERSION}-1_amd64.deb" "${BUILD_DIR}/chroot/tmp/"

    chroot "${BUILD_DIR}/chroot" bash -c "
        export DEBIAN_FRONTEND=noninteractive
        apt-get update
        dpkg -i /tmp/stumpfworks-nas_${VERSION}-1_amd64.deb || apt-get install -f -y
        rm /tmp/stumpfworks-nas_${VERSION}-1_amd64.deb
    "
else
    # Option 2: Build and install from source
    echo "Building StumpfWorks NAS from source..."

    # Install build dependencies
    chroot "${BUILD_DIR}/chroot" bash -c "
        export DEBIAN_FRONTEND=noninteractive
        apt-get update
        apt-get install -y golang-go nodejs npm git
    "

    # Copy source to chroot
    mkdir -p "${BUILD_DIR}/chroot/tmp/stumpfworks-build"
    cp -r "${PROJECT_ROOT}/backend" "${BUILD_DIR}/chroot/tmp/stumpfworks-build/"
    cp -r "${PROJECT_ROOT}/frontend" "${BUILD_DIR}/chroot/tmp/stumpfworks-build/"
    cp "${PROJECT_ROOT}/config.example.yaml" "${BUILD_DIR}/chroot/tmp/stumpfworks-build/"

    # Build in chroot
    chroot "${BUILD_DIR}/chroot" bash -c "
        cd /tmp/stumpfworks-build/backend
        go build -o /usr/bin/stumpfworks-server ./cmd/stumpfworks-server

        cd /tmp/stumpfworks-build/frontend
        npm ci
        npm run build

        mkdir -p /usr/share/stumpfworks-nas
        cp -r dist/* /usr/share/stumpfworks-nas/

        mkdir -p /etc/stumpfworks-nas
        cp /tmp/stumpfworks-build/config.example.yaml /etc/stumpfworks-nas/config.yaml

        mkdir -p /var/lib/stumpfworks-nas
    "

    # Install systemd service
    cp "${PROJECT_ROOT}/debian/stumpfworks-nas.service" "${BUILD_DIR}/chroot/lib/systemd/system/"

    # Enable service
    chroot "${BUILD_DIR}/chroot" systemctl enable stumpfworks-nas.service

    # Cleanup build files
    rm -rf "${BUILD_DIR}/chroot/tmp/stumpfworks-build"
fi

# Install StumpfWorks NAS dependencies
chroot "${BUILD_DIR}/chroot" bash -c "
    export DEBIAN_FRONTEND=noninteractive
    apt-get install -y \
        samba \
        smbclient \
        smartmontools \
        nfs-kernel-server \
        lvm2 \
        mdadm \
        docker.io \
        btrfs-progs
    apt-get clean
"

# Unmount filesystems
umount "${BUILD_DIR}/chroot/dev/pts"
umount "${BUILD_DIR}/chroot/dev"
umount "${BUILD_DIR}/chroot/sys"
umount "${BUILD_DIR}/chroot/proc"

echo "✓ StumpfWorks NAS installed"
