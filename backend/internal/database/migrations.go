// Revision: 2025-11-17 | Author: Claude | Version: 1.1.2
package database

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

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

	// Create default admin user if none exists
	if err := createDefaultAdmin(); err != nil {
		return err
	}

	return nil
}

// createDefaultAdmin creates a default admin user if no admin exists
func createDefaultAdmin() error {
	var count int64
	if err := DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		// Generate secure random password
		password, err := generateSecurePassword(16)
		if err != nil {
			return fmt.Errorf("failed to generate admin password: %w", err)
		}

		admin := &models.User{
			Username: "admin",
			Email:    "admin@stumpfworks.local",
			FullName: "System Administrator",
			Role:     "admin",
			IsActive: true,
		}

		if err := admin.SetPassword(password); err != nil {
			return err
		}

		if err := DB.Create(admin).Error; err != nil {
			return err
		}

		// Print to STDOUT (not logs!) so it's visible but not persisted
		fmt.Fprintln(os.Stdout, "\n"+separator(80))
		fmt.Fprintln(os.Stdout, "🔐 DEFAULT ADMIN USER CREATED")
		fmt.Fprintln(os.Stdout, separator(80))
		fmt.Fprintf(os.Stdout, "   Username: %s\n", admin.Username)
		fmt.Fprintf(os.Stdout, "   Password: %s\n", password)
		fmt.Fprintln(os.Stdout, separator(80))
		fmt.Fprintln(os.Stdout, "⚠️  IMPORTANT:")
		fmt.Fprintln(os.Stdout, "   - Save this password NOW! It will not be shown again.")
		fmt.Fprintln(os.Stdout, "   - Change this password immediately after first login!")
		fmt.Fprintln(os.Stdout, "   - This password is NOT stored in logs for security.")
		fmt.Fprintln(os.Stdout, separator(80)+"\n")

		logger.Info("Default admin user created with random password (displayed on console)")
	}

	return nil
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
