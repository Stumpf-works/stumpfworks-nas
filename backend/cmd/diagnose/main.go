// Diagnose tool for password reset issues
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
)

func main() {
	fmt.Println("=== Password Reset System Diagnostics ===\n")

	// Load config
	configPath := os.Getenv("STUMPFWORKS_CONFIG")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		cfg, _ = config.Load("")
	}

	// Initialize logger
	if err := logger.InitLogger("info", false); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	fmt.Println("✅ Database connected\n")

	// Check if password_reset_tokens table exists
	var tableExists bool
	err = database.DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name='password_reset_tokens'").Row().Scan(&tableExists)
	if err != nil {
		fmt.Printf("✅ password_reset_tokens table: EXISTS\n")
	} else {
		fmt.Printf("❌ password_reset_tokens table: NOT FOUND\n")
		fmt.Println("   Run migrations to create the table\n")
	}

	// Check admin user
	var admin models.User
	if err := database.DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		fmt.Printf("❌ Admin user not found: %v\n\n", err)
	} else {
		fmt.Printf("✅ Admin user found:\n")
		fmt.Printf("   ID: %d\n", admin.ID)
		fmt.Printf("   Username: %s\n", admin.Username)
		fmt.Printf("   Email: %s\n", admin.Email)
		fmt.Printf("   Role: %s\n", admin.Role)
		fmt.Printf("   Active: %v\n", admin.IsActive)
		if admin.LastLoginAt != nil {
			fmt.Printf("   Last Login: %s\n", admin.LastLoginAt.Format(time.RFC3339))
		} else {
			fmt.Printf("   Last Login: Never\n")
		}
		fmt.Println()
	}

	// Check for existing reset tokens
	var tokenCount int64
	database.DB.Model(&models.PasswordResetToken{}).Count(&tokenCount)
	fmt.Printf("📊 Password reset tokens in database: %d\n", tokenCount)

	if tokenCount > 0 {
		var tokens []models.PasswordResetToken
		database.DB.Preload("User").Order("created_at DESC").Limit(5).Find(&tokens)
		fmt.Println("\n📋 Recent tokens (last 5):")
		for i, token := range tokens {
			fmt.Printf("\n   %d. Token: %s...\n", i+1, token.Token[:16])
			fmt.Printf("      User: %s (ID: %d)\n", token.User.Username, token.UserID)
			fmt.Printf("      Created: %s\n", token.CreatedAt.Format(time.RFC3339))
			fmt.Printf("      Expires: %s\n", token.ExpiresAt.Format(time.RFC3339))
			if token.UsedAt != nil {
				fmt.Printf("      Used: %s ✅\n", token.UsedAt.Format(time.RFC3339))
			} else if token.IsExpired() {
				fmt.Printf("      Status: EXPIRED ⏰\n")
			} else {
				fmt.Printf("      Status: VALID ✅\n")
			}
		}
	}

	fmt.Println("\n=== Diagnostics Complete ===")
}
