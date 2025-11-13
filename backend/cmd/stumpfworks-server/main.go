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
	"github.com/Stumpf-works/stumpfworks-nas/internal/docker"
	"github.com/Stumpf-works/stumpfworks-nas/internal/plugins"
	"github.com/Stumpf-works/stumpfworks-nas/internal/updates"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
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

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()

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
