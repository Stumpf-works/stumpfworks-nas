// Revision: 2025-11-23 | Author: Claude | Version: 1.2.0
package database

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// RunMigrations runs all database migrations
func RunMigrations() error {
	logger.Info("Running database migrations...")

	// Auto-migrate models
	if err := DB.AutoMigrate(
		&models.User{},
		&models.UserGroup{},
		&models.Share{},
		&models.DiskLabel{},
		&models.AuditLog{},
		&models.FailedLoginAttempt{},
		&models.IPBlock{},
		&models.AlertConfig{},
		&models.AlertLog{},
		&models.ScheduledTask{},
		&models.TaskExecution{},
		&models.TwoFactorAuth{},
		&models.TwoFactorBackupCode{},
		&models.TwoFactorAttempt{},
		&models.SystemMetric{},
		&models.HealthScore{},
		&models.MonitoringConfig{},
		&models.AddonInstallation{},
		// Add more models here as they are created
	); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")

	// Add performance indexes
	if err := AddPerformanceIndexes(); err != nil {
		logger.Warn("Failed to add performance indexes (non-fatal)", zap.Error(err))
	}

	// NOTE: Default admin user creation removed.
	// Users must now use the Setup Wizard on first access to create the initial admin account.

	return nil
}

// AddPerformanceIndexes adds database indexes for improved query performance
func AddPerformanceIndexes() error {
	logger.Info("Adding performance indexes...")

	// Index for user username lookups (login queries)
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return err
	}

	// Indexes for audit log queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(created_at DESC)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)").Error; err != nil {
		return err
	}
	// Composite index for filtered time-based queries
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_user_timestamp ON audit_logs(user_id, created_at DESC)").Error; err != nil {
		return err
	}

	// Index for share name lookups
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_shares_name ON shares(name)").Error; err != nil {
		return err
	}

	// Index for failed login attempts by IP
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_failed_logins_ip ON failed_login_attempts(ip_address)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_failed_logins_timestamp ON failed_login_attempts(attempted_at DESC)").Error; err != nil {
		return err
	}

	// Index for scheduled tasks
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_enabled ON scheduled_tasks(enabled)").Error; err != nil {
		return err
	}

	// Index for system metrics timestamp
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON system_metrics(timestamp DESC)").Error; err != nil {
		return err
	}

	logger.Info("Performance indexes added successfully")
	return nil
}
