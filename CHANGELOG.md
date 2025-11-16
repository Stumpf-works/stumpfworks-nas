# Changelog

All notable changes to Stumpf.Works NAS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### 🚀 Advanced Features & Plugin System

#### Added

**Phase 4: Plugin System with Runtime Execution**
- Plugin Runtime Engine for sandboxed execution
- `StartPlugin()` / `StopPlugin()` / `RestartPlugin()` - Full lifecycle control
- `GetPluginStatus()` - Runtime status monitoring
- `ListRunningPlugins()` - Active plugin tracking
- Plugin process isolation with separate execution contexts
- Automatic stdout/stderr capture with plugin-tagged logging
- Graceful shutdown handling (5s timeout, force kill if needed)
- Crash detection and automatic status updates
- Environment variables: PLUGIN_ID, PLUGIN_DIR, NAS_API_URL
- Configuration loading from plugin.json
- Example "Hello World" plugin in Go
- Comprehensive Plugin SDK documentation (200+ lines)
- Multi-language support (Go, Python, Node.js examples)
- 5 new API endpoints for runtime control

### 🚀 Advanced Sharing Protocols & Critical Improvements

#### Added

**Phase 2: Advanced iSCSI Management** (19 new methods)
- `AddACL()` - ACL management for initiator access control
- `RemoveACL()` - Remove initiator from target ACL
- `SetCHAPAuth()` - Configure CHAP authentication for secure access
- `DisableCHAPAuth()` - Disable CHAP authentication
- `AddPortal()` - Network portal management for multi-path
- `RemovePortal()` - Remove network portal
- `CreateBlockBackstore()` - Block device backstore creation
- `CreateFileIOBackstore()` - FileIO backstore for file-backed LUNs
- `DeleteBackstore()` - Remove backstores
- `AddLUN()` - LUN management and mapping
- `RemoveLUN()` - Remove LUN from target
- `EnableTarget()` / `DisableTarget()` - Target state control
- `SetTargetAttribute()` - Configure target parameters
- `GetSessions()` - Active session monitoring
- `SaveConfig()` / `RestoreConfig()` / `ClearConfig()` - Configuration management
- Enhanced `ListTargets()` with proper targetcli output parsing

**Phase 2: Advanced WebDAV Management** (10 new methods)
- `EnableApacheWebDAV()` - Enable WebDAV modules in Apache
- `CreateApacheVHost()` - Virtual host creation for WebDAV shares
- `DeleteApacheVHost()` - Remove virtual host configuration
- `CreateHTPasswdFile()` - HTPasswd-based authentication
- `AddUserToHTPasswd()` / `RemoveUserFromHTPasswd()` - User management
- `SetPermissions()` - File system permission control
- `EnableSSL()` - SSL/TLS configuration for secure WebDAV
- `TestConfiguration()` - Apache configuration validation
- `GetApacheLogs()` - Log retrieval for debugging

**Phase 2: Advanced FTP Management** (20 new methods)
- `UpdateVsftpdConfig()` - Dynamic vsftpd configuration management
- `ConfigureVsftpdTLS()` - FTPS/TLS setup with modern security
- `EnableAnonymousFTP()` / `EnableLocalUsers()` / `EnableWriteAccess()` - Access control
- `SetBandwidthLimit()` - Upload/download rate limiting
- `SetMaxClients()` / `SetMaxPerIP()` - Connection limiting
- `SetPasvPorts()` - Passive mode port configuration
- `EnableChrootJail()` - Security isolation for FTP users
- `AddVirtualUser()` - Virtual user management
- `BanIP()` / `UnbanIP()` - IP-based access control
- `EnableLogging()` / `GetLogs()` - Comprehensive logging
- `GetActiveConnections()` - Connection monitoring
- `BackupConfiguration()` / `RestoreConfiguration()` - Config management
- `TestConfiguration()` - Configuration validation

**File Service Improvements**
- Proper owner/group extraction using syscall for Unix file metadata
- `getUsername()` and `getGroupname()` helper functions with UID/GID lookup
- File info now displays actual Unix user/group names instead of generic "system"
- Graceful fallback to "uid:N" / "gid:N" format when lookup fails

**Share Management Security**
- ValidGroups validation in `CreateShare()` and `UpdateShare()`
- System group validation using `os/user.LookupGroup()`
- Prevents invalid group names from being added to share configurations
- Consistent with existing ValidUsers validation pattern

**WebSocket Enhancements**
- Full subscription/unsubscription logic with channel tracking
- Subscriptions map in Client struct for per-channel filtering
- Confirmation messages on subscribe/unsubscribe events
- `IsSubscribed()` and `GetSubscriptions()` helper methods
- Enables real-time channel-based message filtering

**WebSocket Security**
- Proper origin checking in WebSocket upgrader
- Allows localhost, 127.0.0.1, and same-origin requests
- Logs and denies unknown origins for security
- Production-ready with configurable whitelist support

**User API Validation**
- Comprehensive CreateUser validation:
  - Username: required, 3-100 characters
  - Email: required format validation
  - Password: required, minimum 8 characters
  - Role: required, must be admin/user/guest
- UpdateUser validation for optional fields:
  - Password: minimum 8 characters if provided
  - Role: must be admin/user/guest if provided
- Improved error messages for better developer/user experience

#### Fixed

**iSCSI Target Management**
- Implemented proper targetcli output parsing in `ListTargets()`
- Extracts IQN from targetcli tree structure
- Added `extractTargetName()` helper for friendly name parsing
- Returns proper `ISCSITarget` structs instead of empty list

**Code Quality**
- ✅ Zero TODOs remaining in backend codebase
- ✅ All critical security validations implemented
- ✅ Proper error handling across all new features
- ✅ Comprehensive logging for debugging

#### Technical Details
- All implementations use centralized ShellExecutor pattern
- No breaking API changes
- Zero new dependencies added
- Maintains backward compatibility
- Build Status: ✅ All code compiles successfully

---

## [1.1.0] - 2025-11-16

### 🎉 Major System Architecture Improvements

This release introduces the **StumpfWorks System Library v1.1.0**, a complete architectural refactoring that centralizes all system management operations into a unified, thread-safe API.

#### Added

**System Library v1.1.0**
- Centralized SystemLibrary providing unified API for all system operations
- Thread-safe operation with proper mutex locking for all subsystems
- Global instance management with Initialize() and Get() methods
- Comprehensive health monitoring with HealthCheck() returning detailed subsystem status
- Graceful degradation when optional components are unavailable

**Storage Management Enhancements**
- Added `GetPool(name string)` method to ZFSManager for retrieving specific ZFS pools
- Added `GetArray(name string)` method to RAIDManager for retrieving specific RAID arrays
- Improved error handling and type safety across storage operations
- Better integration with the centralized System Library

**Shell Executor Refactoring**
- Centralized ShellExecutor in `internal/system/executor` package
- Eliminated duplicate executor implementations across subsystems
- Improved security with dry-run support for testing
- Better error messages and logging
- Consistent command execution across all system components

**Network Management**
- Complete interface management with bonding support
- Firewall configuration and rule management
- DNS configuration management
- Fixed type references for BondConfig (proper Slaves field mapping)

**Health & Monitoring**
- Enhanced health check system returning detailed subsystem status
- HealthStatus struct with timestamp, overall status, and per-subsystem health
- SubsystemHealth tracking for: Storage, Network, Sharing, Users, Shell, Metrics
- Proper error handling and reporting in health checks

#### Fixed

**Backend Compilation & Type Safety**
- Fixed all API handler compilation errors in syslib.go
- Corrected type references: `sharing.SambaShare`, `sharing.NFSExport`, `network.BondConfig`
- Fixed field access: Network.Interface → Network.Interfaces
- Fixed BondConfig initialization: Interfaces → Slaves field
- Renamed Health() to HealthCheck() for consistency across codebase
- Updated all handlers to use proper error handling with HealthCheck()
- Added required package imports (sharing, network) to API handlers
- Resolved SystemMetrics type conflicts and LoadAvg import issues

**Code Organization**
- Eliminated circular dependencies in system packages
- Centralized executor interface and implementation
- Improved package structure and separation of concerns
- Better type safety throughout the codebase

#### Changed

**Architecture**
- Migrated from scattered system operations to centralized System Library
- All subsystems now accessed through single SystemLibrary instance
- Improved initialization order and dependency management
- Better error propagation from subsystems to API layer

**API Handlers**
- Updated all handlers to use centralized System Library
- Improved error responses with proper error types
- Better validation of requests and responses
- Enhanced logging for debugging and monitoring

**Documentation**
- Updated README.md to v1.1.0 with comprehensive improvements
- Added "System Library Components" section detailing all subsystems
- Enhanced architecture diagram showing System Library integration
- Added "What Makes StumpfWorks NAS Different?" section
- Updated metrics: 170 features, 160+ endpoints, 22 API handlers
- Improved installation instructions with APT repository support
- Updated feature completion tracking (Phase 5 complete)

#### Technical Details

**System Library Components:**
- **Storage Manager**: ZFS pools, RAID arrays, disk operations, SMART monitoring
- **Network Manager**: Interfaces, bonding, firewall rules, DNS configuration
- **Sharing Manager**: Samba (SMB) and NFS exports with user permissions
- **User Manager**: System users, authentication, permissions
- **Metrics Collector**: Real-time system metrics and health monitoring
- **Shell Executor**: Secure command execution with dry-run support

**Build & Quality:**
- Backend compiles without errors ✅
- Frontend builds successfully ✅
- All type mismatches resolved ✅
- Comprehensive error handling ✅
- Thread-safe operations ✅

#### Migration Notes

For users upgrading from v1.0.0:
- No breaking API changes - all existing endpoints remain compatible
- System Library initialization is automatic on startup
- Improved error messages may change log output format
- Health check response format enhanced with more detailed information

---

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

