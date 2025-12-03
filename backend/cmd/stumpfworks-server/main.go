// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/ad"
	"github.com/Stumpf-works/stumpfworks-nas/internal/addons"
	"github.com/Stumpf-works/stumpfworks-nas/internal/alerts"
	"github.com/Stumpf-works/stumpfworks-nas/internal/alertrules"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api/handlers"
	"github.com/Stumpf-works/stumpfworks-nas/internal/audit"
	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/internal/backup"
	"github.com/Stumpf-works/stumpfworks-nas/internal/cloudbackup"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/dependencies"
	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/internal/metrics"
	"github.com/Stumpf-works/stumpfworks-nas/internal/network"
	"github.com/Stumpf-works/stumpfworks-nas/internal/plugins"
	"github.com/Stumpf-works/stumpfworks-nas/internal/scheduler"
	"github.com/Stumpf-works/stumpfworks-nas/internal/storage"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/filesystem"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/ha"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/lxc"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/vm"
	"github.com/Stumpf-works/stumpfworks-nas/internal/system/vpn"
	"github.com/Stumpf-works/stumpfworks-nas/internal/twofa"
	"github.com/Stumpf-works/stumpfworks-nas/internal/ups"
	"github.com/Stumpf-works/stumpfworks-nas/internal/updates"
	"github.com/Stumpf-works/stumpfworks-nas/internal/usergroups"
	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
	"go.uber.org/zap"
)

const (
	AppName    = "Stumpf.Works NAS"
	AppVersion = "1.0.0"
)

func main() {
	// Parse command line flags
	resetAdminPassword := flag.String("reset-admin-password", "", "Reset password for admin user (provide username)")
	flag.Parse()

	fmt.Printf("%s v%s\n", AppName, AppVersion)

	// Load configuration
	configPath := os.Getenv("STUMPFWORKS_CONFIG")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		// If config file doesn't exist, use defaults
		cfg, _ = config.Load("")
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Logging.Level, cfg.IsDevelopment()); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Handle password reset command
	if *resetAdminPassword != "" {
		handlePasswordReset(cfg, *resetAdminPassword)
		return
	}

	fmt.Println("Starting server...")

	logger.Info("Configuration loaded",
		zap.String("environment", cfg.App.Environment),
		zap.String("version", cfg.App.Version))

	// Perform system health check
	performSystemHealthCheck(cfg)

	// Check system dependencies
	if cfg.Dependencies.CheckOnStartup {
		if err := checkDependencies(cfg); err != nil {
			logger.Warn("Dependency check failed - some features may not work",
				zap.Error(err))
			// Don't fail startup - system can work with missing optional packages
		}
	} else {
		logger.Info("Dependency check disabled in configuration")
	}

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()

	// Initialize System Library
	if err := system.Initialize(nil); err != nil {
		logger.Fatal("Failed to initialize system library", zap.Error(err))
	}
	logger.Info("System library initialized")

	// Start System Library background tasks (metrics collection, etc.)
	if err := system.MustGet().Start(); err != nil {
		logger.Warn("Failed to start system library background tasks",
			zap.Error(err),
			zap.String("message", "Metrics collection may be limited"))
	}

	// Initialize Samba user manager (non-fatal if Samba not installed)
	if err := initializeSambaUserManager(); err != nil {
		logger.Warn("Samba user manager initialization failed",
			zap.Error(err),
			zap.String("message", "Samba user sync disabled - users will only work for web access"))
	} else {
		logger.Info("Samba user manager initialized")
	}

	// Initialize Unix group manager (non-fatal if commands not available)
	if err := initializeUnixGroupManager(); err != nil {
		logger.Warn("Unix group manager initialization failed",
			zap.Error(err),
			zap.String("message", "Unix group sync disabled - groups will only work in database"))
	} else {
		logger.Info("Unix group manager initialized")
	}

	// Ensure default shares exist (creates default shares on first run)
	if err := storage.EnsureDefaultShares(); err != nil {
		logger.Warn("Failed to ensure default shares",
			zap.Error(err),
			zap.String("message", "You may need to create shares manually"))
	} else {
		logger.Info("Default shares verified")
	}

	// Fix permissions for all existing shares
	if err := storage.FixExistingSharePermissions(); err != nil {
		logger.Warn("Failed to fix share permissions",
			zap.Error(err),
			zap.String("message", "Some shares may have incorrect permissions"))
	} else {
		logger.Info("Share permissions verified and fixed")
	}

	// Repair Samba configuration (fixes common misconfigurations)
	if err := storage.RepairSambaConfig(); err != nil {
		logger.Warn("Failed to repair Samba configuration",
			zap.Error(err),
			zap.String("message", "Samba shares may not work correctly - check /etc/samba/smb.conf"))
	} else {
		logger.Info("Samba configuration verified and repaired if needed")
	}

	// Mount all persisted volumes from database (ensures volumes persist across reboots)
	if err := storage.MountPersistedVolumes(); err != nil {
		logger.Warn("Failed to mount some volumes from database",
			zap.Error(err),
			zap.String("message", "Some volumes may be offline - check storage page for details"))
	} else {
		logger.Info("All persisted volumes mounted successfully")
	}

	// Restore all network bridges from database (ensures bridges persist across reboots)
	if err := restoreNetworkBridges(); err != nil {
		logger.Warn("Failed to restore some network bridges from database",
			zap.Error(err),
			zap.String("message", "Some network bridges may be offline - check network page for details"))
	} else {
		logger.Info("All persisted network bridges restored successfully")
	}

	// Initialize file service
	if err := handlers.InitFileService(); err != nil {
		logger.Fatal("Failed to initialize file service", zap.Error(err))
	}
	logger.Info("File service initialized")

	// Initialize ACL service (non-fatal if ACL tools not available)
	if err := initializeACL(); err != nil {
		logger.Warn("ACL service initialization failed",
			zap.Error(err),
			zap.String("message", "ACL features will be disabled"))
	} else {
		logger.Info("ACL service initialized")
	}

	// Initialize Quota service (non-fatal if quota tools not available)
	if err := initializeQuota(); err != nil {
		logger.Warn("Quota service initialization failed",
			zap.Error(err),
			zap.String("message", "Quota features will be disabled"))
	} else {
		logger.Info("Quota service initialized")
	}

	// Initialize DRBD service (non-fatal if DRBD tools not available)
	if err := initializeDRBD(); err != nil {
		logger.Warn("DRBD service initialization failed",
			zap.Error(err),
			zap.String("message", "DRBD features will be disabled"))
	} else {
		logger.Info("DRBD service initialized")
	}

	// Initialize Pacemaker/Corosync service (non-fatal if not available)
	if err := initializePacemaker(); err != nil {
		logger.Warn("Pacemaker service initialization failed",
			zap.Error(err),
			zap.String("message", "Pacemaker/Corosync features will be disabled"))
	} else {
		logger.Info("Pacemaker/Corosync service initialized")
	}

	// Initialize Keepalived service (non-fatal if not available)
	if err := initializeKeepalived(); err != nil {
		logger.Warn("Keepalived service initialization failed",
			zap.Error(err),
			zap.String("message", "Virtual IP (Keepalived) features will be disabled"))
	} else {
		logger.Info("Keepalived service initialized")
	}

	// Initialize Addon Manager (always enabled)
	initializeAddonManager()

	// Initialize VM Manager (non-fatal, requires VM Manager addon)
	if err := initializeVMManager(); err != nil {
		logger.Warn("VM Manager initialization failed",
			zap.Error(err),
			zap.String("message", "VM management features will be disabled. Install VM Manager addon to enable."))
	} else {
		logger.Info("VM Manager initialized")
	}

	// Initialize LXC Manager (non-fatal, requires LXC Manager addon)
	var lxcManagerInstance *lxc.LXCManager
	if err := initializeLXCManager(); err != nil {
		logger.Warn("LXC Manager initialization failed",
			zap.Error(err),
			zap.String("message", "LXC management features will be disabled. Install LXC Manager addon to enable."))
	} else {
		logger.Info("LXC Manager initialized")
		// Get the LXC manager instance for container restoration
		shell := system.MustGet().Shell
		lxcManagerInstance, _ = lxc.NewLXCManager(shell)

		// Restore autostart containers from database
		if err := restoreLXCContainers(lxcManagerInstance); err != nil {
			logger.Warn("Failed to restore some LXC containers from database",
				zap.Error(err),
				zap.String("message", "Some autostart containers may not have started - check LXC page for details"))
		} else {
			logger.Info("All autostart LXC containers restored successfully")
		}
	}

	// Initialize VPN Manager (non-fatal, requires VPN Server addon)
	if err := initializeVPNManager(); err != nil {
		logger.Warn("VPN Manager initialization failed",
			zap.Error(err),
			zap.String("message", "VPN management features will be disabled. Install VPN Server addon to enable."))
	} else {
		logger.Info("VPN Manager initialized")
	}

	// Initialize Docker service (non-fatal if not available)
	if err := initializeDocker(); err != nil {
		logger.Warn("Docker not available",
			zap.Error(err),
			zap.String("message", "Docker features will be disabled"))
	} else {
		logger.Info("Docker service initialized and available")
	}

	// Initialize Plugin service (non-fatal if fails)
	if err := initializePlugins(); err != nil {
		logger.Warn("Plugin service initialization failed",
			zap.Error(err),
			zap.String("message", "Plugin features may be limited"))
	} else {
		logger.Info("Plugin service initialized")
	}

	// Initialize Backup service (non-fatal if fails)
	if err := initializeBackup(); err != nil {
		logger.Warn("Backup service initialization failed",
			zap.Error(err),
			zap.String("message", "Backup features may be limited"))
	} else {
		logger.Info("Backup service initialized")
	}

	// Initialize Cloud Backup service (non-fatal if fails)
	if err := initializeCloudBackup(); err != nil {
		logger.Warn("Cloud backup service initialization failed",
			zap.Error(err),
			zap.String("message", "Cloud backup features may be limited - ensure rclone is installed"))
	} else {
		logger.Info("Cloud backup service initialized")
	}

	// Initialize UPS service (non-fatal if fails)
	if err := initializeUPS(); err != nil {
		logger.Warn("UPS service initialization failed",
			zap.Error(err),
			zap.String("message", "UPS monitoring features may be limited"))
	} else {
		logger.Info("UPS service initialized")
	}

	// Initialize Active Directory service (non-fatal if fails)
	if err := initializeAD(); err != nil {
		logger.Warn("Active Directory service initialization failed",
			zap.Error(err),
			zap.String("message", "AD features will be disabled"))
	} else {
		logger.Info("Active Directory service initialized")
	}

	// Initialize Active Directory Domain Controller service (non-fatal if fails)
	if err := initializeADDC(); err != nil {
		logger.Warn("AD Domain Controller service initialization failed",
			zap.Error(err),
			zap.String("message", "AD DC features will be disabled"))
	} else {
		logger.Info("AD Domain Controller service initialized")
	}

	// Initialize Audit Log service
	if err := initializeAuditLog(); err != nil {
		logger.Warn("Audit log service initialization failed",
			zap.Error(err),
			zap.String("message", "Audit logging may be limited"))
	} else {
		logger.Info("Audit log service initialized")
	}

	// Initialize Failed Login Tracking service
	if err := initializeFailedLoginService(); err != nil {
		logger.Warn("Failed login service initialization failed",
			zap.Error(err),
			zap.String("message", "Failed login tracking may be limited"))
	} else {
		logger.Info("Failed login tracking service initialized")
	}

	// Initialize Update service
	if err := initializeUpdateService(); err != nil {
		logger.Warn("Update service initialization failed",
			zap.Error(err),
			zap.String("message", "Update checking may be limited"))
	} else {
		logger.Info("Update service initialized")
	}

	// Initialize Alert service
	if err := initializeAlertService(); err != nil {
		logger.Warn("Alert service initialization failed",
			zap.Error(err),
			zap.String("message", "Email alerts may be disabled"))
	} else {
		logger.Info("Alert service initialized")
	}

	// Initialize Scheduler service
	if err := initializeScheduler(); err != nil {
		logger.Warn("Scheduler service initialization failed",
			zap.Error(err),
			zap.String("message", "Scheduled tasks may be disabled"))
	} else {
		logger.Info("Scheduler service initialized and started")
	}

	// Initialize Two-Factor Authentication service
	if err := initializeTwoFA(); err != nil {
		logger.Warn("Two-Factor Authentication service initialization failed",
			zap.Error(err),
			zap.String("message", "2FA may be disabled"))
	} else {
		logger.Info("Two-Factor Authentication service initialized")
	}

	// Initialize Metrics service
	if err := initializeMetrics(); err != nil {
		logger.Warn("Metrics service initialization failed",
			zap.Error(err),
			zap.String("message", "Metrics collection may be disabled"))
	} else {
		logger.Info("Metrics service initialized and started")
	}

	// Initialize Alert Rules service
	if err := initializeAlertRules(); err != nil {
		logger.Warn("Alert Rules service initialization failed",
			zap.Error(err),
			zap.String("message", "Custom alert rules may be disabled"))
	} else {
		logger.Info("Alert Rules service initialized and started")
	}

	// Create HTTP router
	router := api.NewRouter(cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting",
			zap.String("address", server.Addr),
			zap.String("environment", cfg.App.Environment))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully",
		zap.String("address", server.Addr),
		zap.String("health", "http://"+server.Addr+"/health"),
		zap.String("api", "http://"+server.Addr+"/api/v1"))

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

// initializeDocker initializes the Docker service
// Returns error if Docker is not available, but this is non-fatal
func initializeDocker() error {
	_, err := docker.Initialize()
	return err
}

// initializePlugins initializes the Plugin service
// Returns error if plugin service fails to initialize, but this is non-fatal
func initializePlugins() error {
	_, err := plugins.Initialize("")
	return err
}

// initializeBackup initializes the Backup service
// Returns error if backup service fails to initialize, but this is non-fatal
func initializeBackup() error {
	_, err := backup.Initialize("")
	return err
}

// initializeCloudBackup initializes the cloud backup service
func initializeCloudBackup() error {
	_, err := cloudbackup.Initialize()
	return err
}

// initializeAD initializes the Active Directory service
// Returns error if AD service fails to initialize, but this is non-fatal
func initializeAD() error {
	// Initialize with default config (disabled by default)
	_, err := ad.Initialize(nil)
	return err
}

// initializeADDC initializes the Active Directory Domain Controller service
// Returns error if AD DC service fails to initialize, but this is non-fatal
func initializeADDC() error {
	_, err := ad.InitializeDC()
	return err
}

// initializeAuditLog initializes the Audit Log service
// Returns error if audit log service fails to initialize, but this is non-fatal
func initializeAuditLog() error {
	_, err := audit.Initialize()
	return err
}

// initializeFailedLoginService initializes the Failed Login Tracking service
// Returns error if service fails to initialize, but this is non-fatal
func initializeFailedLoginService() error {
	_, err := auth.InitializeFailedLoginService()
	return err
}

// initializeUpdateService initializes the Update service
// Returns error if service fails to initialize, but this is non-fatal
func initializeUpdateService() error {
	_, err := updates.Initialize()
	return err
}

// initializeAlertService initializes the Alert service
// Returns error if service fails to initialize, but this is non-fatal
func initializeAlertService() error {
	_, err := alerts.Initialize()
	return err
}

// initializeUPS initializes the UPS monitoring service
// Returns error if service fails to initialize, but this is non-fatal
func initializeUPS() error {
	_, err := ups.Initialize()
	return err
}

// initializeScheduler initializes the Scheduler service and starts it
// Returns error if service fails to initialize, but this is non-fatal
func initializeScheduler() error {
	service, err := scheduler.Initialize()
	if err != nil {
		return err
	}
	return service.Start()
}

// initializeTwoFA initializes the Two-Factor Authentication service
// Returns error if service fails to initialize, but this is non-fatal
func initializeTwoFA() error {
	_, err := twofa.Initialize()
	return err
}

// initializeSambaUserManager initializes the Samba user synchronization manager
// Returns error if service fails to initialize, but this is non-fatal
func initializeSambaUserManager() error {
	manager := users.InitSambaUserManager()
	if !manager.IsEnabled() {
		return fmt.Errorf("samba not installed")
	}
	return nil
}

// initializeUnixGroupManager initializes the Unix group synchronization manager
// Returns error if service fails to initialize, but this is non-fatal
func initializeUnixGroupManager() error {
	manager := usergroups.InitUnixGroupManager()
	if !manager.IsEnabled() {
		return fmt.Errorf("required Unix group commands not available")
	}
	return nil
}

// initializeMetrics initializes the Metrics collection service and starts it
// Returns error if service fails to initialize, but this is non-fatal
func initializeMetrics() error {
	service, err := metrics.Initialize()
	if err != nil {
		return err
	}
	return service.Start()
}

// initializeAlertRules initializes the Alert Rules evaluation service and starts it
// Returns error if service fails to initialize, but this is non-fatal
func initializeAlertRules() error {
	service, err := alertrules.Initialize()
	if err != nil {
		return err
	}
	return service.Start()
}

// initializeACL initializes the ACL (Access Control List) service
// Returns error if ACL tools are not installed, but this is non-fatal
func initializeACL() error {
	shell := system.MustGet().Shell
	aclManager, err := filesystem.NewACLManager(shell)
	if err != nil {
		return err
	}
	handlers.InitACLManager(aclManager)
	return nil
}

// initializeQuota initializes the Disk Quota service
// Returns error if quota tools are not installed, but this is non-fatal
func initializeQuota() error {
	shell := system.MustGet().Shell
	quotaManager, err := filesystem.NewQuotaManager(shell)
	if err != nil {
		return err
	}
	handlers.InitQuotaManager(quotaManager)
	return nil
}

// initializeDRBD initializes the DRBD (High Availability) service
// Returns error if DRBD tools are not installed, but this is non-fatal
func initializeDRBD() error {
	shell := system.MustGet().Shell
	drbdManager, err := ha.NewDRBDManager(shell)
	if err != nil {
		return err
	}
	handlers.InitDRBDManager(drbdManager)
	return nil
}

// initializePacemaker initializes the Pacemaker/Corosync (Cluster HA) service
// Returns error if Pacemaker tools are not installed, but this is non-fatal
func initializePacemaker() error {
	shell := system.MustGet().Shell
	pacemakerManager, err := ha.NewPacemakerManager(shell)
	if err != nil {
		return err
	}
	handlers.InitPacemakerManager(pacemakerManager)
	return nil
}

// initializeKeepalived initializes the Keepalived (VIP Management) service
// Returns error if Keepalived is not installed, but this is non-fatal
func initializeKeepalived() error {
	shell := system.MustGet().Shell
	keepalivedManager, err := ha.NewKeepalivedManager(shell)
	if err != nil {
		return err
	}
	handlers.InitKeepalivedManager(keepalivedManager)
	return nil
}

// initializeAddonManager initializes the Addon Manager
// This is always enabled and manages installable addons
func initializeAddonManager() {
	shell := system.MustGet().Shell
	addonManager := addons.NewManager(shell)
	handlers.InitAddonManager(addonManager)
	logger.Info("Addon manager initialized")
}

// initializeVMManager initializes the VM Manager
// Returns error if libvirt is not installed, but this is non-fatal
func initializeVMManager() error {
	shell := system.MustGet().Shell
	vmManager, err := vm.NewLibvirtManager(shell)
	if err != nil {
		return err
	}
	handlers.InitVMManager(vmManager)
	return nil
}

// initializeLXCManager initializes the LXC Manager
// Returns error if LXC is not installed, but this is non-fatal
func initializeLXCManager() error {
	shell := system.MustGet().Shell
	lxcManager, err := lxc.NewLXCManager(shell)
	if err != nil {
		return err
	}
	handlers.InitLXCManager(lxcManager)
	return nil
}

// initializeVPNManager initializes the VPN Manager
// Returns error if initialization fails, but this is non-fatal
func initializeVPNManager() error {
	shell := system.MustGet().Shell
	vpnManager := vpn.NewVPNManager(shell)
	handlers.InitVPNManager(vpnManager)
	return nil
}

// checkDependencies checks and optionally installs system dependencies
func checkDependencies(cfg *config.Config) error {
	logger.Info("Checking system dependencies",
		zap.String("mode", cfg.Dependencies.InstallMode))

	// Create installer with configured mode
	mode := dependencies.CheckOnly
	switch cfg.Dependencies.InstallMode {
	case "auto":
		mode = dependencies.AutoInstall
	case "interactive":
		mode = dependencies.Interactive
	case "check":
		mode = dependencies.CheckOnly
	default:
		logger.Warn("Unknown install mode, using 'check'",
			zap.String("mode", cfg.Dependencies.InstallMode))
	}

	installer := dependencies.NewInstaller(mode)
	return installer.CheckAndInstall()
}

// performSystemHealthCheck runs a comprehensive system health check
func performSystemHealthCheck(cfg *config.Config) {
	logger.Info("Running system health check...")

	report := sysutil.PerformSystemHealthCheck()

	// Log summary
	logger.Info("System health check completed",
		zap.String("overallStatus", report.OverallStatus),
		zap.Int("totalChecks", report.Summary.TotalChecks),
		zap.Int("passed", report.Summary.Passed),
		zap.Int("warnings", report.Summary.Warnings),
		zap.Int("errors", report.Summary.Errors),
		zap.Int("missing", report.Summary.Missing),
		zap.Int("requiredMissing", report.Summary.RequiredMissing))

	// In development mode, print full report
	if cfg.IsDevelopment() {
		fmt.Println()
		report.PrintReport()
		fmt.Println()

		// Also save JSON report to file for debugging
		if jsonReport, err := report.ToJSON(); err == nil {
			jsonFile := "./health-check.json"
			if err := os.WriteFile(jsonFile, []byte(jsonReport), 0644); err == nil {
				logger.Info("Health check report saved",
					zap.String("file", jsonFile))
			}
		}
	}

	// Log warnings for missing optional components
	for _, check := range report.Checks {
		if check.Status == "warning" || check.Status == "missing" {
			logger.Warn("Optional component not available",
				zap.String("component", check.Name),
				zap.String("status", check.Status),
				zap.String("message", check.Message))
		}
	}

	// Fail startup if required components are missing
	if report.Summary.RequiredMissing > 0 {
		logger.Error("Required system components are missing - cannot start server")
		for _, check := range report.Checks {
			if check.Required && (check.Status == "error" || check.Status == "missing") {
				logger.Error("Missing required component",
					zap.String("component", check.Name),
					zap.String("message", check.Message))
			}
		}
		os.Exit(1)
	}
}

// handlePasswordReset handles the password reset command for admin users
func handlePasswordReset(cfg *config.Config, username string) {
	fmt.Println("\n" + separator(80))
	fmt.Println("🔐 PASSWORD RESET UTILITY")
	fmt.Println(separator(80))

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	db := database.GetDB()
	if db == nil {
		fmt.Println("❌ Database connection failed")
		os.Exit(1)
	}

	// Find user by username
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Printf("❌ User '%s' not found in database\n", username)
		fmt.Println(separator(80))
		os.Exit(1)
	}

	// Verify user is admin
	if user.Role != "admin" {
		fmt.Printf("❌ User '%s' is not an administrator (role: %s)\n", username, user.Role)
		fmt.Println("⚠️  Password reset is only available for admin users")
		fmt.Println(separator(80))
		os.Exit(1)
	}

	// Generate new secure password
	newPassword, err := generateSecurePassword(16)
	if err != nil {
		fmt.Printf("❌ Failed to generate password: %v\n", err)
		os.Exit(1)
	}

	// Set new password
	if err := user.SetPassword(newPassword); err != nil {
		fmt.Printf("❌ Failed to set password: %v\n", err)
		os.Exit(1)
	}

	// Update user in database
	if err := db.Save(&user).Error; err != nil {
		fmt.Printf("❌ Failed to update user in database: %v\n", err)
		os.Exit(1)
	}

	// Success - display new password
	fmt.Printf("✅ Password reset successful for admin user: %s\n", user.Username)
	fmt.Println(separator(80))
	fmt.Printf("   Username: %s\n", user.Username)
	fmt.Printf("   Email:    %s\n", user.Email)
	fmt.Printf("   New Password: %s\n", newPassword)
	fmt.Println(separator(80))
	fmt.Println("⚠️  IMPORTANT SECURITY NOTES:")
	fmt.Println("   - This password will NOT be shown again!")
	fmt.Println("   - Save this password in a secure location NOW")
	fmt.Println("   - Change this password after logging in")
	fmt.Println("   - This password is NOT stored in any logs")
	fmt.Println(separator(80))
	fmt.Println()

	logger.Info("Admin password reset completed",
		zap.String("username", user.Username),
		zap.String("user_id", fmt.Sprintf("%d", user.ID)))
}

// generateSecurePassword generates a cryptographically secure random password
func generateSecurePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Use base64 encoding for readable password
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// separator creates a visual separator line
func separator(width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s += "="
	}
	return s
}

// restoreNetworkBridges restores all network bridges from database on startup
// This ensures bridges persist across reboots
func restoreNetworkBridges() error {
	return network.RestoreAllBridges()
}

// restoreLXCContainers restores all autostart LXC containers from database on startup
// This ensures containers with autostart enabled are started after reboot
func restoreLXCContainers(lxcManager *lxc.LXCManager) error {
	if lxcManager == nil {
		logger.Info("LXC manager not initialized, skipping container restoration")
		return nil
	}
	return lxcManager.RestoreAllContainers()
}
