# StumpfWorks NAS OS - ISO Builder

This directory contains the configuration for building a bootable ISO image of StumpfWorks NAS OS, similar to TrueNAS SCALE, Unraid, or Proxmox.

## Overview

StumpfWorks NAS OS is a custom Debian-based operating system that includes:
- Minimal Debian base system
- StumpfWorks NAS pre-installed
- Automated installation wizard
- First-boot configuration UI
- All required dependencies

## Prerequisites

### On Debian/Ubuntu Build Machine

```bash
sudo apt update
sudo apt install -y \
    debootstrap \
    squashfs-tools \
    xorriso \
    isolinux \
    syslinux-efi \
    grub-pc-bin \
    grub-efi-amd64-bin \
    mtools
```

### Live-Build Tool (Optional, for advanced customization)

```bash
sudo apt install -y live-build
```

## Building the ISO

### Quick Build (Automated Script)

```bash
cd iso-builder
sudo ./build-iso.sh
```

The ISO will be created in `iso-builder/output/stumpfworks-nas-1.0.0-amd64.iso`

### Manual Build (Step by Step)

```bash
# 1. Prepare build environment
cd iso-builder
sudo ./scripts/prepare-environment.sh

# 2. Build base system
sudo ./scripts/build-base-system.sh

# 3. Install StumpfWorks NAS
sudo ./scripts/install-stumpfworks.sh

# 4. Configure system
sudo ./scripts/configure-system.sh

# 5. Create ISO
sudo ./scripts/create-iso.sh
```

## ISO Features

- **Bootable USB/DVD**: Works on both BIOS and UEFI systems
- **Live Mode**: Try before installing
- **Automated Installer**: Simple installation wizard
- **First-Boot Setup**: Web-based configuration on first boot
- **Pre-configured**: All dependencies and services ready
- **Minimal Size**: ~1.5GB ISO

## Installation Process

1. Boot from USB/DVD
2. Choose "Install StumpfWorks NAS OS"
3. Select installation disk
4. Configure network (optional)
5. System installs automatically
6. Reboot and access web UI at http://YOUR_IP:8080

## Customization

### Adding Packages

Edit `config/package-list.txt` to add additional packages.

### Preseed Configuration

Edit `config/preseed.cfg` to customize automated installation.

### Branding

- Logo: `branding/logo.png`
- Splash: `branding/splash.png`
- Colors: `branding/theme.conf`

## Boot Menu Options

1. **Install StumpfWorks NAS OS** - Standard installation
2. **Install (Expert Mode)** - Advanced installation options
3. **Live Mode** - Try without installing
4. **Rescue Mode** - System recovery tools
5. **Memory Test** - Test RAM

## Directory Structure

```
iso-builder/
├── README.md                    # This file
├── build-iso.sh                 # Main build script
├── config/                      # Configuration files
│   ├── package-list.txt         # Packages to include
│   ├── preseed.cfg              # Automated installation config
│   ├── grub.cfg                 # GRUB bootloader config
│   └── isolinux.cfg             # Legacy boot config
├── scripts/                     # Build scripts
│   ├── prepare-environment.sh
│   ├── build-base-system.sh
│   ├── install-stumpfworks.sh
│   ├── configure-system.sh
│   └── create-iso.sh
├── branding/                    # Visual customization
│   ├── logo.png
│   ├── splash.png
│   └── theme.conf
├── hooks/                       # Live-build hooks
│   └── 9999-customize.hook.chroot
└── output/                      # Build output (created)
    └── stumpfworks-nas-VERSION-amd64.iso
```

## Testing the ISO

### Using QEMU/KVM

```bash
# Install QEMU
sudo apt install qemu-system-x86

# Test ISO
qemu-system-x86_64 \
    -cdrom output/stumpfworks-nas-1.0.0-amd64.iso \
    -boot d \
    -m 4096 \
    -enable-kvm
```

### Using VirtualBox

1. Create new VM (Linux, Debian 64-bit)
2. Allocate 4GB RAM minimum
3. Attach ISO to optical drive
4. Boot VM

### Creating Bootable USB

```bash
# Find USB device
lsblk

# Write ISO to USB (replace /dev/sdX with your USB device)
sudo dd if=output/stumpfworks-nas-1.0.0-amd64.iso of=/dev/sdX bs=4M status=progress && sync
```

## Troubleshooting

### Build fails with "Permission denied"

Run build script with sudo:
```bash
sudo ./build-iso.sh
```

### ISO doesn't boot in UEFI mode

Ensure grub-efi packages are installed:
```bash
sudo apt install grub-efi-amd64-bin
```

### ISO is too large

Reduce package list in `config/package-list.txt`

## Release Process

1. Build ISO with release script:
   ```bash
   sudo ./build-iso.sh --release
   ```

2. Test ISO on real hardware

3. Upload to GitHub releases:
   ```bash
   gh release upload v1.0.0 output/stumpfworks-nas-1.0.0-amd64.iso
   ```

## Requirements

- **Build Machine**: Debian 11+ or Ubuntu 20.04+ with 10GB free space
- **Target Hardware**: x86-64 CPU, 2GB RAM minimum, 20GB disk
- **Network**: Internet connection for package downloads

## More Information

- Debian Live Manual: https://live-team.pages.debian.net/live-manual/
- Debian Preseed: https://www.debian.org/releases/stable/amd64/apb.html
- GRUB Manual: https://www.gnu.org/software/grub/manual/

## Support

- GitHub Issues: https://github.com/Stumpf-works/stumpfworks-nas/issues
- Documentation: https://github.com/Stumpf-works/stumpfworks-nas/wiki
