package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/config"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/internal/api"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/internal/core"
	"github.com/stumpfworks/stumpfworks-nas/plugins/vpn-server/pkg/database"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = "1.0.0"
	buildTime  = "unknown"
)

func main() {
	flag.Parse()

	// Print banner
	printBanner()

	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Connect to database
	log.Println("Connecting to database...")
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create VPN manager
	log.Println("Initializing VPN manager...")
	vpnManager, err := core.NewVPNManager(cfg, db)
	if err != nil {
		log.Fatalf("Failed to create VPN manager: %v", err)
	}
	defer vpnManager.Close()

	// Start VPN servers
	log.Println("Starting VPN servers...")
	ctx := context.Background()
	if err := vpnManager.Start(ctx); err != nil {
		log.Printf("Warning: Failed to start VPN servers: %v", err)
	}

	// Setup HTTP router
	log.Println("Setting up HTTP server...")
	router := api.SetupRouter(cfg, vpnManager)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Stop VPN servers
	log.Println("Stopping VPN servers...")
	if err := vpnManager.Stop(); err != nil {
		log.Printf("Error stopping VPN servers: %v", err)
	}

	log.Println("Server exited successfully")
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Configure logger
	logLevel := logger.Silent
	if cfg.General.LogLevel == "debug" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ███████╗████████╗██╗   ██╗███╗   ███╗██████╗ ███████╗  ║
║   ██╔════╝╚══██╔══╝██║   ██║████╗ ████║██╔══██╗██╔════╝  ║
║   ███████╗   ██║   ██║   ██║██╔████╔██║██████╔╝█████╗    ║
║   ╚════██║   ██║   ██║   ██║██║╚██╔╝██║██╔═══╝ ██╔══╝    ║
║   ███████║   ██║   ╚██████╔╝██║ ╚═╝ ██║██║     ██║       ║
║   ╚══════╝   ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝       ║
║                                                           ║
║              VPN Server - Multi-Protocol                 ║
║                                                           ║
║   WireGuard | OpenVPN | PPTP | L2TP/IPsec               ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
	Version: %s
	Build: %s
`
	fmt.Printf(banner, version, buildTime)
	fmt.Println()
}
