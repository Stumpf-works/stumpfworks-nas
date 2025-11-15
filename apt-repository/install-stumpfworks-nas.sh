#!/bin/bash
# StumpfWorks NAS - One-Line Installer
# Quick installation script for users

set -e

# Configuration
REPO_URL="https://stumpf-works.github.io/stumpfworks-apt-repo"
GPG_KEY_URL="${REPO_URL}/stumpfworks.gpg"
KEYRING_PATH="/usr/share/keyrings/stumpfworks-archive-keyring.gpg"
SOURCES_LIST="/etc/apt/sources.list.d/stumpfworks-nas.list"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    cat << "EOF"

╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║            StumpfWorks NAS - Quick Installer                 ║
║                                                              ║
║  Production-ready NAS management system                      ║
║  Alternative to TrueNAS, Unraid, and Synology DSM           ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

EOF
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "This script must be run as root"
        echo ""
        echo "Please run with sudo:"
        echo "  curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/install.sh | sudo bash"
        exit 1
    fi
}

check_system() {
    log_info "Checking system compatibility..."

    # Check if Debian-based
    if [ ! -f /etc/debian_version ]; then
        log_error "This script only works on Debian-based systems"
        exit 1
    fi

    # Check architecture
    ARCH=$(dpkg --print-architecture)
    if [ "$ARCH" != "amd64" ]; then
        log_warn "Only amd64 architecture is currently supported"
        log_warn "Your architecture: $ARCH"
    fi

    # Get OS info
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        log_info "Detected: $PRETTY_NAME"
    fi

    log_info "✓ System is compatible"
}

install_dependencies() {
    log_info "Checking dependencies..."

    local deps_needed=false

    if ! command -v curl &> /dev/null; then
        deps_needed=true
    fi

    if ! command -v gpg &> /dev/null; then
        deps_needed=true
    fi

    if [ "$deps_needed" = true ]; then
        log_info "Installing required dependencies..."
        apt-get update -qq
        apt-get install -y curl gnupg >/dev/null 2>&1
        log_info "✓ Dependencies installed"
    else
        log_info "✓ All dependencies satisfied"
    fi
}

add_gpg_key() {
    log_info "Adding StumpfWorks GPG key..."

    # Create keyrings directory if not exists
    mkdir -p /usr/share/keyrings

    # Download and add GPG key
    if curl -fsSL "$GPG_KEY_URL" | gpg --dearmor -o "$KEYRING_PATH" 2>/dev/null; then
        log_info "✓ GPG key added"
    else
        log_error "Failed to download GPG key"
        exit 1
    fi
}

add_repository() {
    log_info "Adding StumpfWorks APT repository..."

    # Create sources.list.d directory if not exists
    mkdir -p /etc/apt/sources.list.d

    # Add repository
    echo "deb [signed-by=${KEYRING_PATH}] ${REPO_URL} stable main" > "$SOURCES_LIST"

    log_info "✓ Repository added"
}

update_package_lists() {
    log_info "Updating package lists..."

    if apt-get update -qq 2>&1 | grep -q "stumpfworks"; then
        log_info "✓ Package lists updated"
    else
        log_warn "Repository may not be accessible yet"
    fi
}

install_stumpfworks() {
    log_info "Installing StumpfWorks NAS..."

    if apt-get install -y stumpfworks-nas >/dev/null 2>&1; then
        log_info "✓ StumpfWorks NAS installed successfully!"
    else
        log_error "Installation failed"
        log_info "Try manually: sudo apt install stumpfworks-nas"
        exit 1
    fi
}

print_success() {
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                              ║${NC}"
    echo -e "${GREEN}║  ✓ StumpfWorks NAS installed successfully!                   ║${NC}"
    echo -e "${GREEN}║                                                              ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    # Get server IP
    SERVER_IP=$(ip -4 addr show | grep -oP '(?<=inet\s)\d+(\.\d+){3}' | grep -v '127.0.0.1' | head -n1)

    echo "Next steps:"
    echo ""
    echo "  1. Configure the system (optional):"
    echo "     sudo nano /etc/stumpfworks-nas/config.yaml"
    echo ""
    echo "  2. Start the service:"
    echo "     sudo systemctl start stumpfworks-nas"
    echo ""
    echo "  3. Enable auto-start on boot:"
    echo "     sudo systemctl enable stumpfworks-nas"
    echo ""
    echo "  4. Check service status:"
    echo "     sudo systemctl status stumpfworks-nas"
    echo ""
    echo "  5. Access the web interface:"

    if [ -n "$SERVER_IP" ]; then
        echo "     http://${SERVER_IP}:8080"
    else
        echo "     http://YOUR_SERVER_IP:8080"
    fi

    echo ""
    echo "  Default credentials:"
    echo "     Username: admin"
    echo "     Password: admin"
    echo "     (⚠️  Change immediately after first login!)"
    echo ""
    echo "  View logs:"
    echo "     sudo journalctl -u stumpfworks-nas -f"
    echo ""
    echo "  Update in the future:"
    echo "     sudo apt update && sudo apt upgrade"
    echo ""
    echo "Documentation: https://github.com/Stumpf-works/stumpfworks-nas"
    echo ""
}

main() {
    print_banner
    check_root
    check_system
    install_dependencies
    add_gpg_key
    add_repository
    update_package_lists
    install_stumpfworks
    print_success
}

main "$@"
