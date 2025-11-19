// Set password tool - directly set user password with proper bcrypt hashing
package main

import (
	"fmt"
	"os"

	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database"
	"github.com/Stumpf-works/stumpfworks-nas/internal/database/models"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: set-password <username> <new-password>")
		fmt.Println("Example: set-password admin MyNewPassword123")
		os.Exit(1)
	}

	username := os.Args[1]
	newPassword := os.Args[2]

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
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Find user
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Printf("❌ User '%s' not found: %v\n", username, err)
		os.Exit(1)
	}

	fmt.Printf("Found user: %s (ID: %d, Email: %s)\n", user.Username, user.ID, user.Email)
	fmt.Println("Setting new password...")

	// Set new password using the same method as the User model
	if err := user.SetPassword(newPassword); err != nil {
		fmt.Printf("❌ Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// Save to database
	if err := database.DB.Save(&user).Error; err != nil {
		fmt.Printf("❌ Failed to save user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Password successfully updated for user '%s'\n", username)
	fmt.Println("You can now log in with the new password.")
}
