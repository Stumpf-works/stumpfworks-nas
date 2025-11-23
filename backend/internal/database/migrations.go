// Revision: 2025-11-23 | Author: Claude | Version: 1.2.0
package database

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
)

// RunMigrations runs all database migrations
func RunMigrations() error {
	logger.Info("Running database migrations...")

	// Auto-migrate models
	if err := DB.AutoMigrate(
		&models.User{},
		&models.UserGroup{},
		&models.Share{},
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
		// Add more models here as they are created
	); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")

	// NOTE: Default admin user creation removed.
	// Users must now use the Setup Wizard on first access to create the initial admin account.

	return nil
}
