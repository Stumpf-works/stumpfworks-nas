#!/bin/bash
# StumpfWorks NAS - APT Repository Initialization Script
# This script sets up the APT repository structure

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="${SCRIPT_DIR}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
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
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║     StumpfWorks NAS - APT Repository Setup               ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
EOF
    echo ""
}

check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()
    local required_deps=("dpkg-scanpackages" "gpg" "gzip")

    for dep in "${required_deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing_deps+=("$dep")
        fi
    done

    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Install with: sudo apt install dpkg-dev gnupg gzip"
        exit 1
    fi

    log_info "✓ All dependencies satisfied"
}

create_directory_structure() {
    log_info "Creating directory structure..."

    mkdir -p "${REPO_ROOT}/pool/main"
    mkdir -p "${REPO_ROOT}/dists/stable/main/binary-amd64"
    mkdir -p "${REPO_ROOT}/dists/stable/main/binary-arm64"

    log_info "✓ Directory structure created"
}

check_gpg_key() {
    log_info "Checking for GPG key..."

    if gpg --list-keys "contact@stumpf.works" &>/dev/null; then
        log_info "✓ GPG key found for contact@stumpf.works"
        return 0
    else
        log_warn "No GPG key found for contact@stumpf.works"
        log_info ""
        log_info "To generate a GPG key, run:"
        log_info "  gpg --full-generate-key"
        log_info ""
        log_info "  Choose:"
        log_info "    - Key type: RSA and RSA"
        log_info "    - Key size: 4096"
        log_info "    - Expiration: 2y (2 years)"
        log_info "    - Name: Stumpf.Works Team"
        log_info "    - Email: contact@stumpf.works"
        log_info ""
        log_warn "Repository will be created but NOT signed!"
        log_warn "You can sign it later with: ./sign-repository.sh"
        return 1
    fi
}

export_public_key() {
    if gpg --list-keys "contact@stumpf.works" &>/dev/null; then
        log_info "Exporting public GPG key..."
        gpg --armor --export contact@stumpf.works > "${REPO_ROOT}/stumpfworks.gpg"
        log_info "✓ Public key exported to stumpfworks.gpg"
    fi
}

add_packages() {
    log_info "Looking for .deb packages..."

    local deb_count=0

    # Look for .deb files in parent directory
    if ls "${SCRIPT_DIR}"/../*.deb 1> /dev/null 2>&1; then
        cp "${SCRIPT_DIR}"/../*.deb "${REPO_ROOT}/pool/main/"
        deb_count=$(ls "${REPO_ROOT}/pool/main/"*.deb 2>/dev/null | wc -l)
        log_info "✓ Copied ${deb_count} .deb package(s) from parent directory"
    else
        log_warn "No .deb packages found in parent directory"
        log_info "Build the package first with: dpkg-buildpackage -b -uc -us"
    fi
}

generate_packages_file() {
    log_info "Generating Packages index..."

    cd "${REPO_ROOT}"

    # Generate for amd64
    dpkg-scanpackages --arch amd64 pool/ > dists/stable/main/binary-amd64/Packages 2>/dev/null || true
    gzip -k -f dists/stable/main/binary-amd64/Packages

    log_info "✓ Packages file generated"
}

generate_release_file() {
    log_info "Generating Release file..."

    cd "${REPO_ROOT}/dists/stable"

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
    apt-ftparchive release . >> Release 2>/dev/null || {
        # Fallback if apt-ftparchive not available
        log_warn "apt-ftparchive not available, using basic Release file"
    }

    log_info "✓ Release file generated"
}

sign_repository() {
    if gpg --list-keys "contact@stumpf.works" &>/dev/null; then
        log_info "Signing repository..."

        cd "${REPO_ROOT}/dists/stable"

        # Remove old signatures
        rm -f Release.gpg InRelease

        # Create detached signature
        gpg --default-key contact@stumpf.works -abs -o Release.gpg Release 2>/dev/null

        # Create clearsign signature
        gpg --default-key contact@stumpf.works --clearsign -o InRelease Release 2>/dev/null

        log_info "✓ Repository signed"
    else
        log_warn "Skipping signing (no GPG key)"
    fi
}

create_index_html() {
    log_info "Creating index.html..."

    cat > "${REPO_ROOT}/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>StumpfWorks NAS - APT Repository</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            line-height: 1.6;
            padding: 20px;
        }
        .container {
            max-width: 800px;
            margin: 40px auto;
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            padding: 40px;
        }
        h1 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 2.5em;
        }
        h2 {
            color: #555;
            margin-top: 30px;
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 2px solid #667eea;
        }
        .subtitle {
            color: #888;
            margin-bottom: 30px;
            font-size: 1.1em;
        }
        code {
            background: #f5f5f5;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            color: #d63384;
        }
        pre {
            background: #282c34;
            color: #abb2bf;
            padding: 20px;
            border-radius: 5px;
            overflow-x: auto;
            margin: 15px 0;
        }
        pre code {
            background: none;
            color: inherit;
            padding: 0;
        }
        .button {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 12px 24px;
            border-radius: 5px;
            text-decoration: none;
            margin: 10px 10px 10px 0;
            transition: background 0.3s;
        }
        .button:hover {
            background: #5568d3;
        }
        .info-box {
            background: #e7f3ff;
            border-left: 4px solid #2196F3;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .warning-box {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 StumpfWorks NAS</h1>
        <p class="subtitle">Official APT Repository</p>

        <p>Welcome to the StumpfWorks NAS APT repository. This repository provides easy installation and automatic updates for StumpfWorks NAS on Debian-based systems.</p>

        <h2>📦 Installation</h2>

        <p><strong>Step 1:</strong> Add the GPG key</p>
        <pre><code>curl -fsSL https://stumpf-works.github.io/stumpfworks-apt-repo/stumpfworks.gpg | \
sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks-archive-keyring.gpg</code></pre>

        <p><strong>Step 2:</strong> Add the repository</p>
        <pre><code>echo "deb [signed-by=/usr/share/keyrings/stumpfworks-archive-keyring.gpg] https://stumpf-works.github.io/stumpfworks-apt-repo stable main" | \
sudo tee /etc/apt/sources.list.d/stumpfworks-nas.list</code></pre>

        <p><strong>Step 3:</strong> Install StumpfWorks NAS</p>
        <pre><code>sudo apt update
sudo apt install stumpfworks-nas</code></pre>

        <h2>🔄 Updates</h2>

        <p>Once installed, you'll receive automatic updates through your system's package manager:</p>
        <pre><code>sudo apt update
sudo apt upgrade</code></pre>

        <div class="info-box">
            <strong>💡 Tip:</strong> You can also update directly from the StumpfWorks NAS web interface!
        </div>

        <h2>📚 Documentation</h2>

        <a href="https://github.com/Stumpf-works/stumpfworks-nas" class="button">GitHub Repository</a>
        <a href="https://github.com/Stumpf-works/stumpfworks-nas/wiki" class="button">Documentation</a>
        <a href="https://github.com/Stumpf-works/stumpfworks-nas/releases" class="button">Releases</a>

        <h2>🔐 GPG Key Verification</h2>

        <p>Verify the repository GPG key fingerprint:</p>
        <pre><code>gpg --show-keys /usr/share/keyrings/stumpfworks-archive-keyring.gpg</code></pre>

        <h2>ℹ️ System Requirements</h2>

        <ul style="margin-left: 20px; margin-top: 10px;">
            <li>Debian 11 (Bullseye) or newer</li>
            <li>Ubuntu 20.04 LTS or newer</li>
            <li>2 GB RAM minimum (4 GB recommended)</li>
            <li>20 GB disk space</li>
            <li>x86-64 architecture</li>
        </ul>

        <h2>🆘 Support</h2>

        <p>Need help? Visit our <a href="https://github.com/Stumpf-works/stumpfworks-nas/issues">GitHub Issues</a> or check the <a href="https://github.com/Stumpf-works/stumpfworks-nas/wiki">documentation</a>.</p>

        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #ddd; color: #888; text-align: center;">
            <p>&copy; 2025 Stumpf.Works Team | <a href="https://github.com/Stumpf-works">GitHub</a></p>
        </div>
    </div>
</body>
</html>
EOF

    log_info "✓ index.html created"
}

print_summary() {
    echo ""
    log_info "════════════════════════════════════════════════════════"
    log_info "✓ APT Repository initialized successfully!"
    log_info "════════════════════════════════════════════════════════"
    echo ""
    log_info "Repository location: ${REPO_ROOT}"
    log_info "Public key: ${REPO_ROOT}/stumpfworks.gpg"
    echo ""
    log_info "Next steps:"
    echo ""
    log_info "  1. Review the repository structure"
    log_info "  2. Add .deb packages to pool/main/"
    log_info "  3. Update repository: ./update-repository.sh"
    log_info "  4. Deploy to GitHub Pages (see README.md)"
    echo ""
    log_info "To add this repository to a system:"
    echo ""
    echo "  curl -fsSL https://YOUR-DOMAIN/stumpfworks.gpg | \\"
    echo "    sudo gpg --dearmor -o /usr/share/keyrings/stumpfworks-archive-keyring.gpg"
    echo ""
    echo "  echo \"deb [signed-by=/usr/share/keyrings/stumpfworks-archive-keyring.gpg] https://YOUR-DOMAIN stable main\" | \\"
    echo "    sudo tee /etc/apt/sources.list.d/stumpfworks-nas.list"
    echo ""
    echo "  sudo apt update"
    echo "  sudo apt install stumpfworks-nas"
    echo ""
}

main() {
    print_banner
    check_dependencies
    create_directory_structure
    check_gpg_key
    export_public_key
    add_packages
    generate_packages_file
    generate_release_file
    sign_repository
    create_index_html
    print_summary
}

main "$@"
