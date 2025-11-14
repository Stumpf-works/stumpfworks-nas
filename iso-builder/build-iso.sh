#!/bin/bash
# StumpfWorks NAS OS - ISO Builder
# Builds a bootable Debian-based ISO with StumpfWorks NAS pre-installed

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="1.0.0"
ISO_NAME="stumpfworks-nas-${VERSION}-amd64.iso"
BUILD_DIR="${SCRIPT_DIR}/build"
OUTPUT_DIR="${SCRIPT_DIR}/output"
DEBIAN_VERSION="bookworm"  # Debian 12

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()
    local required_deps=(
        "debootstrap"
        "mksquashfs"
        "xorriso"
        "isolinux"
        "grub-mkrescue"
    )

    for dep in "${required_deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_deps+=("$dep")
        fi
    done

    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Install with: sudo apt install debootstrap squashfs-tools xorriso isolinux grub-pc-bin grub-efi-amd64-bin"
        exit 1
    fi

    log_info "All dependencies satisfied"
}

clean_build() {
    log_info "Cleaning previous build..."
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"
    mkdir -p "${OUTPUT_DIR}"
}

print_banner() {
    cat << "EOF"
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║         StumpfWorks NAS OS - ISO Builder                 ║
║                                                          ║
║   Building Debian-based bootable ISO with NAS system    ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
EOF
    echo ""
}

main() {
    print_banner

    log_info "Starting ISO build for StumpfWorks NAS v${VERSION}"
    log_info "Debian base: ${DEBIAN_VERSION}"
    log_info "Build directory: ${BUILD_DIR}"
    log_info "Output: ${OUTPUT_DIR}/${ISO_NAME}"
    echo ""

    # Checks
    check_root
    check_dependencies

    # Build steps
    clean_build

    log_info "Step 1/5: Preparing build environment..."
    bash "${SCRIPT_DIR}/scripts/prepare-environment.sh" "${BUILD_DIR}" "${DEBIAN_VERSION}"

    log_info "Step 2/5: Building base Debian system..."
    bash "${SCRIPT_DIR}/scripts/build-base-system.sh" "${BUILD_DIR}" "${DEBIAN_VERSION}"

    log_info "Step 3/5: Installing StumpfWorks NAS..."
    bash "${SCRIPT_DIR}/scripts/install-stumpfworks.sh" "${BUILD_DIR}" "${VERSION}"

    log_info "Step 4/5: Configuring system..."
    bash "${SCRIPT_DIR}/scripts/configure-system.sh" "${BUILD_DIR}"

    log_info "Step 5/5: Creating ISO image..."
    bash "${SCRIPT_DIR}/scripts/create-iso.sh" "${BUILD_DIR}" "${OUTPUT_DIR}" "${ISO_NAME}"

    echo ""
    log_info "═══════════════════════════════════════════════════════════"
    log_info "✓ ISO build completed successfully!"
    log_info "═══════════════════════════════════════════════════════════"
    log_info "ISO Location: ${OUTPUT_DIR}/${ISO_NAME}"
    log_info "ISO Size: $(du -h "${OUTPUT_DIR}/${ISO_NAME}" | cut -f1)"
    log_info ""
    log_info "Next steps:"
    log_info "  1. Test ISO: qemu-system-x86_64 -cdrom ${OUTPUT_DIR}/${ISO_NAME} -m 4096 -boot d"
    log_info "  2. Create bootable USB: sudo dd if=${OUTPUT_DIR}/${ISO_NAME} of=/dev/sdX bs=4M status=progress"
    log_info "  3. Upload to GitHub: gh release upload v${VERSION} ${OUTPUT_DIR}/${ISO_NAME}"
    echo ""
}

# Run main function
main "$@"
