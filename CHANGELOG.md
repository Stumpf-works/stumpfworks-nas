# Changelog

All notable changes to Stumpf.Works NAS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-14

### 🎉 Initial Production Release

First production-ready release of Stumpf.Works NAS - a complete bare-metal NAS management system.

#### Added

**Core System**
- Complete web-based NAS management interface with macOS-inspired UI
- Real-time system monitoring with metrics collection (CPU, Memory, Disk, Network)
- System health checks and dependency management
- Comprehensive audit logging for security and compliance

**Storage Management**
- SMB/CIFS share management with Samba integration
- Disk management with SMART monitoring
- LVM volume management
- File browser with upload/download capabilities
- Permission management (owner, group, mode)
- Archive creation/extraction (ZIP, TAR)

**Container Management**
- Full Docker container orchestration
- Docker image management (pull, remove, search)
- Docker volume management
- Docker network management with stability improvements
- Docker Compose stack management
- Container stats and log viewing

**User Management**
- User and group management with Unix/Samba synchronization
- Role-based access control (Admin/User)
- Two-factor authentication (2FA) with TOTP
- Session management with JWT tokens
- Failed login tracking and IP blocking

**Network & Services**
- Network interface configuration
- DNS configuration
- Firewall management (iptables/ufw)
- Network diagnostics (ping, traceroute, netstat)
- Wake-on-LAN support

**Backup & Recovery**
- Backup job scheduling
- Snapshot management
- Backup history tracking

**Additional Features**
- Active Directory integration
- Plugin system for extensibility
- Scheduled task management with cron expressions
- Email and webhook alerts
- System update checking

#### Fixed
- Added ErrorBoundary component to prevent complete UI crashes
- Implemented username/group name resolution in file permissions (no more UID/GID-only display)
- Enhanced Docker Networks stability with defensive validation
- Fixed permission checks for chunked file uploads

#### Technical Details

**Backend**
- Go 1.21+ with Chi router
- SQLite database with GORM
- JWT authentication with bcrypt password hashing
- gopsutil for system metrics
- Comprehensive dependency checking on startup

**Frontend**
- React 18 with TypeScript
- TailwindCSS for styling
- Framer Motion for animations
- Zustand for state management
- macOS-inspired design language

**Architecture**
- RESTful API design
- WebSocket support for real-time updates
- Modular service architecture
- Graceful degradation for missing dependencies

#### Known Limitations
- Advanced ACL support planned for 1.1
- User disk quotas planned for 1.1
- Multi-language support planned for future release
- See TODO.md for complete roadmap

---

## Release Notes

This is the first production-ready release of Stumpf.Works NAS. The system is designed to be a complete alternative to TrueNAS, Unraid, and Synology DSM, running directly on bare metal hardware.

### Installation
See INSTALL.md for detailed installation instructions.

### Upgrading
As this is the initial release, no upgrade path exists yet. Future releases will include upgrade documentation.

### Support
For issues, feature requests, and discussions, please visit:
https://github.com/Stumpf-works/stumpfworks-nas/issues

