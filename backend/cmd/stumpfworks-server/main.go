package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/ad"
	"github.com/Stumpf-works/stumpfworks-nas/internal/alerts"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api"
	"github.com/Stumpf-works/stumpfworks-nas/internal/api/handlers"
	"github.com/Stumpf-works/stumpfworks-nas/internal/audit"
	"github.com/Stumpf-works/stumpfworks-nas/internal/auth"
	"github.com/Stumpf-works/stumpfworks-nas/internal/backup"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/dependencies"
	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/internal/metrics"
	"github.com/Stumpf-works/stumpfworks-nas/internal/plugins"
	"github.com/Stumpf-works/stumpfworks-nas/internal/scheduler"
	"github.com/Stumpf-works/stumpfworks-nas/internal/storage"
	"github.com/Stumpf-works/stumpfworks-nas/internal/twofa"
	"github.com/Stumpf-works/stumpfworks-nas/internal/updates"
	"github.com/Stumpf-works/stumpfworks-nas/internal/usergroups"
	"github.com/Stumpf-works/stumpfworks-nas/internal/users"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/sysutil"
	"go.uber.org/zap"
)

const (
	AppName    = "Stumpf.Works NAS"
	AppVersion = "0.2.1-alpha"
)

func main() {
	fmt.Printf("%s v%s\n", AppName, AppVersion)
	fmt.Println("Starting server...")

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

	// Initialize file service
	if err := handlers.InitFileService(); err != nil {
		logger.Fatal("Failed to initialize file service", zap.Error(err))
	}
	logger.Info("File service initialized")

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

	// Initialize Active Directory service (non-fatal if fails)
	if err := initializeAD(); err != nil {
		logger.Warn("Active Directory service initialization failed",
			zap.Error(err),
			zap.String("message", "AD features will be disabled"))
	} else {
		logger.Info("Active Directory service initialized")
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

// initializeAD initializes the Active Directory service
// Returns error if AD service fails to initialize, but this is non-fatal
func initializeAD() error {
	// Initialize with default config (disabled by default)
	_, err := ad.Initialize(nil)
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
