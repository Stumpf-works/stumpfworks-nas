// Check password tool
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
		fmt.Println("Usage: check-password <username> <password>")
		fmt.Println("Example: check-password admin MyPassword123")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]

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

	// Check password
	if user.CheckPassword(password) {
		fmt.Printf("✅ Password is CORRECT for user '%s'\n", username)
		fmt.Printf("   User ID: %d\n", user.ID)
		fmt.Printf("   Email: %s\n", user.Email)
		fmt.Printf("   Role: %s\n", user.Role)
		fmt.Printf("   Active: %v\n", user.IsActive)
	} else {
		fmt.Printf("❌ Password is INCORRECT for user '%s'\n", username)
		fmt.Printf("   User ID: %d\n", user.ID)
		fmt.Printf("   Email: %s\n", user.Email)
		fmt.Printf("   Role: %s\n", user.Role)
	}
}
