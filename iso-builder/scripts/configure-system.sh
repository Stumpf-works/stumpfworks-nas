#!/bin/bash
# Configure the system for StumpfWorks NAS OS

set -e

BUILD_DIR="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Configuring system..."

# Mount filesystems
mount -t proc none "${BUILD_DIR}/chroot/proc"
mount -t sysfs none "${BUILD_DIR}/chroot/sys"
mount -o bind /dev "${BUILD_DIR}/chroot/dev"
mount -o bind /dev/pts "${BUILD_DIR}/chroot/dev/pts"

# Configure hostname
echo "stumpfworks-nas" > "${BUILD_DIR}/chroot/etc/hostname"

# Configure hosts file
cat > "${BUILD_DIR}/chroot/etc/hosts" << EOF
127.0.0.1       localhost
127.0.1.1       stumpfworks-nas

::1             localhost ip6-localhost ip6-loopback
ff02::1         ip6-allnodes
ff02::2         ip6-allrouters
EOF

# Configure network interfaces (using NetworkManager)
cat > "${BUILD_DIR}/chroot/etc/network/interfaces" << EOF
# This file describes the network interfaces available on your system
# and how to activate them. For more information, see interfaces(5).

source /etc/network/interfaces.d/*

# The loopback network interface
auto lo
iface lo inet loopback

# Network interfaces managed by NetworkManager
EOF

# Configure first-boot script
cat > "${BUILD_DIR}/chroot/usr/local/bin/stumpfworks-firstboot.sh" << 'FBEOF'
#!/bin/bash
# StumpfWorks NAS OS First Boot Configuration

set -e

FIRSTBOOT_FLAG="/var/lib/stumpfworks-nas/.firstboot"

if [ -f "$FIRSTBOOT_FLAG" ]; then
    echo "First boot already completed."
    exit 0
fi

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                                                              ║"
echo "║         Welcome to StumpfWorks NAS OS v1.0.0                 ║"
echo "║                                                              ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Performing first-boot configuration..."
echo ""

# Wait for network
echo "Waiting for network connection..."
timeout=30
while [ $timeout -gt 0 ]; do
    if ip route | grep -q default; then
        echo "✓ Network connected"
        break
    fi
    sleep 1
    timeout=$((timeout - 1))
done

# Get IP address
IP_ADDR=$(ip -4 addr show | grep -oP '(?<=inet\s)\d+(\.\d+){3}' | grep -v '127.0.0.1' | head -n1)

if [ -z "$IP_ADDR" ]; then
    IP_ADDR="<unable to detect>"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✓ StumpfWorks NAS OS is ready!"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Access the Web UI:"
echo "  http://${IP_ADDR}:8080"
echo ""
echo "Default credentials:"
echo "  Username: admin"
echo "  Password: admin"
echo "  (Change immediately after first login!)"
echo ""
echo "SSH access:"
echo "  ssh root@${IP_ADDR}"
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Mark first boot as completed
touch "$FIRSTBOOT_FLAG"
FBEOF

chmod +x "${BUILD_DIR}/chroot/usr/local/bin/stumpfworks-firstboot.sh"

# Create systemd service for first boot
cat > "${BUILD_DIR}/chroot/lib/systemd/system/stumpfworks-firstboot.service" << EOF
[Unit]
Description=StumpfWorks NAS First Boot Configuration
After=network-online.target stumpfworks-nas.service
Wants=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/stumpfworks-firstboot.sh
RemainAfterExit=yes
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
EOF

# Enable first boot service
chroot "${BUILD_DIR}/chroot" systemctl enable stumpfworks-firstboot.service

# Configure root password (default: stumpfworks)
echo "Setting default root password..."
echo "root:stumpfworks" | chroot "${BUILD_DIR}/chroot" chpasswd

# Configure SSH to allow root login (for initial setup only)
sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' "${BUILD_DIR}/chroot/etc/ssh/sshd_config"

# Set timezone to UTC
chroot "${BUILD_DIR}/chroot" ln -sf /usr/share/zoneinfo/UTC /etc/localtime

# Configure locale
echo "en_US.UTF-8 UTF-8" >> "${BUILD_DIR}/chroot/etc/locale.gen"
chroot "${BUILD_DIR}/chroot" locale-gen
echo "LANG=en_US.UTF-8" > "${BUILD_DIR}/chroot/etc/default/locale"

# Install GRUB (for installed system)
chroot "${BUILD_DIR}/chroot" bash -c "
    grub-mkconfig -o /boot/grub/grub.cfg || true
"

# Create MOTD (Message of the Day)
cat > "${BUILD_DIR}/chroot/etc/motd" << 'MOTDEOF'

╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║            StumpfWorks NAS OS v1.0.0                         ║
║                                                              ║
║  A production-ready NAS management system                    ║
║  Alternative to TrueNAS, Unraid, and Synology DSM           ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

Access the Web UI:
  http://<your-server-ip>:8080

Default credentials:
  Username: admin
  Password: admin
  (Change immediately!)

System commands:
  systemctl status stumpfworks-nas  - Check service status
  journalctl -u stumpfworks-nas -f  - View logs
  systemctl restart stumpfworks-nas - Restart service

Documentation:
  https://github.com/Stumpf-works/stumpfworks-nas

MOTDEOF

# Unmount filesystems
umount "${BUILD_DIR}/chroot/dev/pts"
umount "${BUILD_DIR}/chroot/dev"
umount "${BUILD_DIR}/chroot/sys"
umount "${BUILD_DIR}/chroot/proc"

echo "✓ System configured"
