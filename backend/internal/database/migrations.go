package database

import (
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/internal/storage"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// RunMigrations runs all database migrations
func RunMigrations() error {
	logger.Info("Running database migrations...")

	// Auto-migrate models
	if err := DB.AutoMigrate(
		&models.User{},
		&storage.ShareModel{},
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
		admin := &models.User{
			Username: "admin",
			Email:    "admin@stumpfworks.local",
			FullName: "System Administrator",
			Role:     "admin",
			IsActive: true,
		}

		if err := admin.SetPassword("admin"); err != nil {
			return err
		}

		if err := DB.Create(admin).Error; err != nil {
			return err
		}

		logger.Info("Default admin user created",
			zap.String("username", admin.Username),
			zap.String("password", "admin (PLEASE CHANGE THIS!)"))
	}

	return nil
}
