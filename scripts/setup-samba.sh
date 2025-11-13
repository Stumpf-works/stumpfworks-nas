#!/bin/bash
# Stumpf.Works NAS - Samba Installation Script
# This script installs and configures Samba for network shares

set -e

echo "======================================"
echo "  Stumpf.Works NAS - Samba Setup"
echo "======================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: This script must be run as root (use sudo)"
    exit 1
fi

echo "[1/5] Installing Samba..."
apt-get update -qq
apt-get install -y samba samba-common-bin

echo ""
echo "[2/5] Creating Samba shares directory..."
mkdir -p /etc/samba/shares.d
chmod 755 /etc/samba/shares.d

echo ""
echo "[3/5] Backing up original smb.conf..."
if [ -f /etc/samba/smb.conf ]; then
    cp /etc/samba/smb.conf /etc/samba/smb.conf.backup-$(date +%Y%m%d-%H%M%S)
fi

echo ""
echo "[4/5] Configuring Samba..."

# Create base smb.conf if it doesn't exist or is minimal
cat > /etc/samba/smb.conf <<'EOF'
#======================= Global Settings =======================

[global]
   workgroup = WORKGROUP
   server string = Stumpf.Works NAS
   server role = standalone server
   log file = /var/log/samba/log.%m
   max log size = 1000
   logging = file
   panic action = /usr/share/samba/panic-action %d

   # Authentication
   security = user
   passdb backend = tdbsam
   obey pam restrictions = yes
   unix password sync = yes
   passwd program = /usr/bin/passwd %u
   passwd chat = *Enter\snew\s*\spassword:* %n\n *Retype\snew\s*\spassword:* %n\n *password\supdated\ssuccessfully* .
   pam password change = yes
   map to guest = bad user

   # Performance optimizations
   socket options = TCP_NODELAY IPTOS_LOWDELAY SO_RCVBUF=131072 SO_SNDBUF=131072
   read raw = yes
   write raw = yes
   max xmit = 65535
   min receivefile size = 16384
   use sendfile = yes
   aio read size = 16384
   aio write size = 16384

   # Character set
   unix charset = UTF-8
   dos charset = CP850

   # Include share configurations from shares.d
   include = /etc/samba/shares.d/*.conf

#======================= Share Definitions =======================
# Individual shares are configured in /etc/samba/shares.d/
# This allows the Stumpf.Works NAS backend to manage shares dynamically
EOF

echo ""
echo "[5/5] Starting Samba services..."

# Check if system uses systemd or sysvinit
if command -v systemctl &> /dev/null; then
    echo "   Using systemd..."
    systemctl enable smbd
    systemctl enable nmbd
    systemctl restart smbd
    systemctl restart nmbd

    echo ""
    echo "Checking Samba status:"
    systemctl status smbd --no-pager -l || true
elif command -v service &> /dev/null; then
    echo "   Using sysvinit..."
    service smbd start
    service nmbd start
else
    echo "   Starting manually..."
    /usr/sbin/smbd -D
    /usr/sbin/nmbd -D
fi

echo ""
echo "======================================"
echo "  Samba Installation Complete!"
echo "======================================"
echo ""
echo "Next Steps:"
echo "1. The Stumpf.Works NAS backend will automatically create share"
echo "   configurations in /etc/samba/shares.d/ when you create shares"
echo "   through the web interface."
echo ""
echo "2. To add a Samba user (for authentication), run:"
echo "   sudo smbpasswd -a <username>"
echo ""
echo "3. The backend will automatically sync Samba users when you"
echo "   create users through the web interface."
echo ""
echo "Configuration files:"
echo "  - Main config: /etc/samba/smb.conf"
echo "  - Shares:      /etc/samba/shares.d/"
echo "  - Logs:        /var/log/samba/"
echo ""
echo "Useful commands:"
echo "  - Test config:     sudo testparm"
echo "  - List shares:     smbclient -L localhost"
echo "  - Reload config:   sudo systemctl reload smbd"
echo "  - View logs:       sudo tail -f /var/log/samba/log.smbd"
echo ""
