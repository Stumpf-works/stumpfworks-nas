#!/bin/bash
# Prepare build environment for ISO creation

set -e

BUILD_DIR="$1"
DEBIAN_VERSION="$2"

echo "Preparing build environment..."

# Create directory structure
mkdir -p "${BUILD_DIR}/chroot"
mkdir -p "${BUILD_DIR}/iso"
mkdir -p "${BUILD_DIR}/iso/boot/grub"
mkdir -p "${BUILD_DIR}/iso/isolinux"
mkdir -p "${BUILD_DIR}/iso/live"

echo "✓ Build directories created"
