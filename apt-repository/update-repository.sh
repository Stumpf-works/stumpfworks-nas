#!/bin/bash
# StumpfWorks NAS - APT Repository Update Script
# Updates the repository with new packages

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_banner() {
    cat << "EOF"
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║     StumpfWorks NAS - Repository Update                 ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
EOF
    echo ""
}

add_package() {
    local package_path="$1"

    if [ ! -f "$package_path" ]; then
        log_warn "Package not found: $package_path"
        return 1
    fi

    log_info "Adding package: $(basename "$package_path")"
    cp "$package_path" "${REPO_ROOT}/pool/main/"
    log_info "✓ Package copied to pool/main/"
}

regenerate_index() {
    log_info "Regenerating package index..."

    cd "${REPO_ROOT}"

    # Generate Packages file for amd64
    dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages 2>/dev/null
    gzip -k -f dists/stable/main/binary-amd64/Packages

    log_info "✓ Package index regenerated"
}

update_release() {
    log_info "Updating Release file..."

    cd "${REPO_ROOT}/dists/stable"

    # Remove old Release file
    rm -f Release

    # Create new Release file
    cat > Release << EOF
Origin: StumpfWorks
Label: StumpfWorks NAS
Suite: stable
Codename: stable
Version: 1.0
Architectures: amd64 arm64
Components: main
Description: StumpfWorks NAS Official APT Repository
Date: $(date -Ru)
EOF

    # Add checksums
    apt-ftparchive release . >> Release 2>/dev/null || true

    log_info "✓ Release file updated"
}

sign_release() {
    if gpg --list-keys "contact@stumpf.works" &>/dev/null; then
        log_info "Signing Release file..."

        cd "${REPO_ROOT}/dists/stable"

        # Remove old signatures
        rm -f Release.gpg InRelease

        # Sign Release file
        gpg --default-key contact@stumpf.works -abs -o Release.gpg Release 2>/dev/null
        gpg --default-key contact@stumpf.works --clearsign -o InRelease Release 2>/dev/null

        log_info "✓ Release file signed"
    else
        log_warn "No GPG key found, skipping signing"
    fi
}

list_packages() {
    log_info "Packages in repository:"
    echo ""

    if ls "${REPO_ROOT}/pool/main/"*.deb 1> /dev/null 2>&1; then
        for pkg in "${REPO_ROOT}/pool/main/"*.deb; do
            dpkg-deb -I "$pkg" | grep -E "Package:|Version:|Architecture:" | sed 's/^/  /'
            echo ""
        done
    else
        log_warn "No packages found in pool/main/"
    fi
}

print_usage() {
    echo "Usage: $0 [OPTIONS] [PACKAGE.deb]"
    echo ""
    echo "Options:"
    echo "  -a, --add PACKAGE.deb    Add a new package to the repository"
    echo "  -l, --list               List all packages in the repository"
    echo "  -u, --update             Update repository index and signatures"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --add ../stumpfworks-nas_1.0.1-1_amd64.deb"
    echo "  $0 --update"
    echo "  $0 --list"
}

main() {
    if [ $# -eq 0 ]; then
        print_banner
        log_info "Updating repository..."
        regenerate_index
        update_release
        sign_release
        list_packages
        log_info "✓ Repository updated successfully!"
        return 0
    fi

    case "$1" in
        -a|--add)
            if [ -z "$2" ]; then
                echo "Error: Package path required"
                print_usage
                exit 1
            fi
            print_banner
            add_package "$2"
            regenerate_index
            update_release
            sign_release
            log_info "✓ Package added and repository updated!"
            ;;
        -l|--list)
            print_banner
            list_packages
            ;;
        -u|--update)
            print_banner
            regenerate_index
            update_release
            sign_release
            log_info "✓ Repository updated successfully!"
            ;;
        -h|--help)
            print_usage
            ;;
        *.deb)
            print_banner
            add_package "$1"
            regenerate_index
            update_release
            sign_release
            log_info "✓ Package added and repository updated!"
            ;;
        *)
            echo "Error: Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
}

main "$@"
