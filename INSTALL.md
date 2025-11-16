# Installation Guide

Stumpf.Works NAS is a bare-metal NAS operating system designed to run directly on your hardware, similar to TrueNAS, Unraid, or Synology DSM.

## 📌 Important Notes

**About sudo and root access:**
- Most commands in this guide require administrator/root privileges
- If your system has `sudo` installed, use the commands marked "With sudo"
- If your system doesn't have `sudo` (minimal Debian installations), use the "Without sudo" commands
- The `useradd` command is located at `/usr/sbin/useradd` (not in the default user PATH)
- If you get `command not found` errors, you likely need to switch to root using `su -`

## System Requirements

### Minimum Requirements
- **CPU**: 64-bit x86 processor (Intel/AMD) or ARM64
- **RAM**: 2 GB minimum, 4 GB recommended
- **Storage**: 
  - 8 GB for system installation
  - Additional storage for data (recommend separate disks)
- **Network**: Ethernet connection

### Recommended Hardware
- **CPU**: Multi-core 64-bit processor
- **RAM**: 8 GB or more
- **Storage**: 
  - 16 GB SSD for system (for better performance)
  - Multiple HDDs/SSDs for storage pools
- **Network**: Gigabit Ethernet

### Software Requirements
- **OS**: Linux (Debian 11+, Ubuntu 20.04+, or similar)
- **Kernel**: Linux 5.x or newer
- **Required packages**: samba, smbclient, smartmontools
- **Optional packages**: docker, nfs-kernel-server, lvm2, mdadm

---

---

## 🚀 Quick Start for Minimal Systems (Without sudo)

If you're on a minimal Debian installation without sudo, here's a complete installation script you can run as root:

```bash
# 1. Switch to root
su -

# 2. Install dependencies
apt-get update
apt-get install -y samba smbclient smartmontools docker.io wget

# 3. Download binary (adjust version and architecture as needed)
cd /tmp
wget https://github.com/Stumpf-works/stumpfworks-nas/releases/latest/download/stumpfworks-nas-linux-amd64
chmod +x stumpfworks-nas-linux-amd64

# 4. Create user and directories
/usr/sbin/useradd -r -s /bin/false -d /opt/stumpfworks stumpfworks
mkdir -p /opt/stumpfworks /var/lib/stumpfworks /etc/stumpfworks

# 5. Install binary
mv stumpfworks-nas-linux-amd64 /usr/local/bin/stumpfworks-nas
chown root:root /usr/local/bin/stumpfworks-nas
chmod 755 /usr/local/bin/stumpfworks-nas

# 6. Create configuration
cat > /etc/stumpfworks/config.yaml << 'EOF'
app:
  name: "Stumpf.Works NAS"
  environment: "production"
  version: "1.0.0"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

database:
  path: "/var/lib/stumpfworks/nas.db"

logging:
  level: "info"

dependencies:
  check_on_startup: true
  install_mode: "check"
EOF

# 7. Create systemd service
cat > /etc/systemd/system/stumpfworks-nas.service << 'EOF'
[Unit]
Description=Stumpf.Works NAS Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
Environment="STUMPFWORKS_CONFIG=/etc/stumpfworks/config.yaml"
ExecStart=/usr/local/bin/stumpfworks-nas
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# 8. Start service
systemctl daemon-reload
systemctl enable stumpfworks-nas
systemctl start stumpfworks-nas

# 9. Check status
systemctl status stumpfworks-nas

# Exit root shell
exit
```

Now access the web interface at `http://<your-server-ip>:8080` with username `admin` and password `admin`.

---

## Installation Methods

### Method 1: Binary Installation (Recommended)

#### 1. Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/Stumpf-works/stumpfworks-nas/releases):

```bash
# For Linux x86_64
wget https://github.com/Stumpf-works/stumpfworks-nas/releases/download/v1.0.0/stumpfworks-nas-linux-amd64

# For Linux ARM64
wget https://github.com/Stumpf-works/stumpfworks-nas/releases/download/v1.0.0/stumpfworks-nas-linux-arm64

# Make executable
chmod +x stumpfworks-nas-linux-amd64
```

#### 2. Install System Dependencies

**Debian/Ubuntu (with sudo):**
```bash
sudo apt update
sudo apt install -y samba smbclient smartmontools docker.io
```

**Debian/Ubuntu (without sudo):**
```bash
su -
apt-get update
apt-get install -y samba smbclient smartmontools docker.io
exit
```

**For additional features (with sudo):**
```bash
sudo apt install -y nfs-kernel-server lvm2 mdadm
```

**For additional features (without sudo):**
```bash
su -
apt-get install -y nfs-kernel-server lvm2 mdadm
exit
```

#### 3. Create System User

**Option A: If you have sudo installed:**
```bash
sudo /usr/sbin/useradd -r -s /bin/false -d /opt/stumpfworks stumpfworks
sudo mkdir -p /opt/stumpfworks
sudo mkdir -p /var/lib/stumpfworks
sudo mkdir -p /etc/stumpfworks
```

**Option B: If sudo is not installed (minimal systems):**
```bash
# Switch to root user
su -

# Then run these commands as root:
/usr/sbin/useradd -r -s /bin/false -d /opt/stumpfworks stumpfworks
mkdir -p /opt/stumpfworks
mkdir -p /var/lib/stumpfworks
mkdir -p /etc/stumpfworks

# Exit root shell when done
exit
```

**Option C: Install sudo first (recommended):**
```bash
# Switch to root
su -

# Install sudo
apt-get update
apt-get install -y sudo

# Add your user to sudo group (replace 'youruser' with your username)
usermod -aG sudo youruser

# Exit and log back in for changes to take effect
exit
```

**Troubleshooting:**
- If `useradd: command not found`: The command is located at `/usr/sbin/useradd` (not in normal user PATH)
- If `sudo: command not found`: Use Option B or C above
- If `permission denied`: You need root access - use `su -` to switch to root

#### 4. Install Binary

**With sudo:**
```bash
sudo mv stumpfworks-nas-linux-amd64 /usr/local/bin/stumpfworks-nas
sudo chown root:root /usr/local/bin/stumpfworks-nas
sudo chmod 755 /usr/local/bin/stumpfworks-nas
```

**Without sudo (as root):**
```bash
su -
mv stumpfworks-nas-linux-amd64 /usr/local/bin/stumpfworks-nas
chown root:root /usr/local/bin/stumpfworks-nas
chmod 755 /usr/local/bin/stumpfworks-nas
exit
```

#### 5. Create Configuration File

**With sudo:**
```bash
sudo tee /etc/stumpfworks/config.yaml > /dev/null << 'YAML'
app:
  name: "Stumpf.Works NAS"
  environment: "production"
  version: "1.0.0"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

database:
  path: "/var/lib/stumpfworks/nas.db"

logging:
  level: "info"

dependencies:
  check_on_startup: true
  install_mode: "check"
YAML
```

**Without sudo (as root):**
```bash
su -
cat > /etc/stumpfworks/config.yaml << 'YAML'
app:
  name: "Stumpf.Works NAS"
  environment: "production"
  version: "1.0.0"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

database:
  path: "/var/lib/stumpfworks/nas.db"

logging:
  level: "info"

dependencies:
  check_on_startup: true
  install_mode: "check"
YAML
exit
```

#### 6. Create Systemd Service

**With sudo:**
```bash
sudo tee /etc/systemd/system/stumpfworks-nas.service > /dev/null << 'SERVICE'
[Unit]
Description=Stumpf.Works NAS Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
Environment="STUMPFWORKS_CONFIG=/etc/stumpfworks/config.yaml"
ExecStart=/usr/local/bin/stumpfworks-nas
Restart=on-failure
RestartSec=5s

# Security hardening (optional but recommended)
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/stumpfworks /etc/stumpfworks /mnt /srv
ProtectKernelTunables=false
ProtectControlGroups=false

[Install]
WantedBy=multi-user.target
SERVICE
```

**Without sudo (as root):**
```bash
su -
cat > /etc/systemd/system/stumpfworks-nas.service << 'SERVICE'
[Unit]
Description=Stumpf.Works NAS Server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
Environment="STUMPFWORKS_CONFIG=/etc/stumpfworks/config.yaml"
ExecStart=/usr/local/bin/stumpfworks-nas
Restart=on-failure
RestartSec=5s

# Security hardening (optional but recommended)
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/stumpfworks /etc/stumpfworks /mnt /srv
ProtectKernelTunables=false
ProtectControlGroups=false

[Install]
WantedBy=multi-user.target
SERVICE
exit
```

#### 7. Start Service

**With sudo:**
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable stumpfworks-nas

# Start service
sudo systemctl start stumpfworks-nas

# Check status
sudo systemctl status stumpfworks-nas
```

**Without sudo (as root):**
```bash
su -

# Reload systemd
systemctl daemon-reload

# Enable service to start on boot
systemctl enable stumpfworks-nas

# Start service
systemctl start stumpfworks-nas

# Check status
systemctl status stumpfworks-nas

exit
```

#### 8. Access Web Interface

Open your browser and navigate to:
```
http://<your-server-ip>:8080
```

**Default credentials:**
- Username: `admin`
- Password: `admin`

⚠️ **IMPORTANT**: Change the default password immediately after first login!

---

### Method 2: Build from Source

#### 1. Prerequisites

```bash
# Install Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js 20+
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

#### 2. Clone Repository

```bash
git clone https://github.com/Stumpf-works/stumpfworks-nas.git
cd stumpfworks-nas
git checkout v1.0.0
```

#### 3. Build Backend

```bash
cd backend
go build -ldflags="-s -w" -o stumpfworks-nas ./cmd/stumpfworks-server
sudo mv stumpfworks-nas /usr/local/bin/
```

#### 4. Build Frontend

```bash
cd ../frontend
npm install
npm run build
# Frontend files will be in dist/ directory
# Backend serves these automatically from embedded assets
```

#### 5. Follow steps 3-8 from Method 1

---

## Post-Installation

### 1. Configure Firewall

```bash
# Allow web interface
sudo ufw allow 8080/tcp

# Allow Samba
sudo ufw allow 445/tcp
sudo ufw allow 139/tcp

# Allow NFS (if using)
sudo ufw allow 2049/tcp
```

### 2. Configure Storage

1. Navigate to **Storage** > **Disks** in web interface
2. Initialize your data disks
3. Create volumes/pools
4. Create shares

### 3. Create Users

1. Navigate to **Users** in web interface
2. Create user accounts
3. Assign to groups
4. Configure share permissions

### 4. Enable Services

Check **System** > **Services** to enable:
- Samba (for Windows/macOS file sharing)
- NFS (for Linux file sharing)
- Docker (for container management)

---

## Upgrading

### Manual Upgrade

1. Stop service:
   ```bash
   sudo systemctl stop stumpfworks-nas
   ```

2. Backup database:
   ```bash
   sudo cp /var/lib/stumpfworks/nas.db /var/lib/stumpfworks/nas.db.backup
   ```

3. Download new binary and replace:
   ```bash
   sudo mv stumpfworks-nas-new /usr/local/bin/stumpfworks-nas
   ```

4. Start service:
   ```bash
   sudo systemctl start stumpfworks-nas
   ```

---

## Troubleshooting

### Service won't start

Check logs:
```bash
sudo journalctl -u stumpfworks-nas -f
```

### Cannot access web interface

1. Check if service is running:
   ```bash
   sudo systemctl status stumpfworks-nas
   ```

2. Check firewall:
   ```bash
   sudo ufw status
   ```

3. Check if port is listening:
   ```bash
   sudo netstat -tlnp | grep 8080
   ```

### Samba shares not working

1. Check if Samba is installed:
   ```bash
   smbd --version
   ```

2. Check Samba configuration:
   ```bash
   sudo testparm
   ```

3. Restart Samba:
   ```bash
   sudo systemctl restart smbd nmbd
   ```

### Docker features not available

Install Docker:
```bash
sudo apt install -y docker.io
sudo systemctl enable --now docker
sudo systemctl restart stumpfworks-nas
```

---

## Uninstallation

```bash
# Stop and disable service
sudo systemctl stop stumpfworks-nas
sudo systemctl disable stumpfworks-nas

# Remove files
sudo rm /usr/local/bin/stumpfworks-nas
sudo rm /etc/systemd/system/stumpfworks-nas.service
sudo rm -rf /etc/stumpfworks
sudo rm -rf /var/lib/stumpfworks

# Reload systemd
sudo systemctl daemon-reload
```

⚠️ **Note**: This does NOT remove your data or shares! Only the NAS management software.

---

## Support

For issues and questions:
- GitHub Issues: https://github.com/Stumpf-works/stumpfworks-nas/issues
- Documentation: https://github.com/Stumpf-works/stumpfworks-nas/tree/main/docs

