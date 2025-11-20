// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize initializes the database connection
func Initialize(cfg *config.Config) error {
	var err error

	// Ensure data directory exists
	dataDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Configure GORM logger
	gormLogLevel := gormlogger.Silent
	if cfg.Logging.Development {
		gormLogLevel = gormlogger.Info
	}

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	switch cfg.Database.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.Database.Path), gormConfig)
	case "postgres", "postgresql":
		// Build PostgreSQL DSN
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Database,
			cfg.Database.SSLMode,
		)
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		return fmt.Errorf("unsupported database driver: %s (supported: sqlite, postgres)", cfg.Database.Driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Use config values or defaults
	maxOpenConns := cfg.Database.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	maxIdleConns := cfg.Database.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}
	connMaxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
	if err != nil || connMaxLifetime == 0 {
		connMaxLifetime = 5 * time.Minute
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	logger.Info("Database connected successfully",
		zap.String("driver", cfg.Database.Driver),
		zap.String("path", cfg.Database.Path))

	// Run migrations
	if err := RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
