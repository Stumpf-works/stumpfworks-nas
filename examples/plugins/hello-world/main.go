package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config represents the plugin configuration
type Config struct {
	Message  string `json:"message"`
	Interval int    `json:"interval"`
}

func main() {
	// Get plugin environment
	pluginID := os.Getenv("PLUGIN_ID")
	pluginDir := os.Getenv("PLUGIN_DIR")

	fmt.Printf("Starting plugin: %s\n", pluginID)
	fmt.Printf("Plugin directory: %s\n", pluginDir)

	// Load configuration from plugin.json
	config, err := loadConfig(pluginDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for periodic messages
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	// Main loop
	fmt.Println("Plugin started successfully!")

	for {
		select {
		case <-ticker.C:
			// Print configured message
			fmt.Println(config.Message)

		case sig := <-sigChan:
			fmt.Printf("Received signal: %v. Shutting down gracefully...\n", sig)
			return
		}
	}
}

// loadConfig loads the plugin configuration
func loadConfig(pluginDir string) (*Config, error) {
	configPath := pluginDir + "/plugin.json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var manifest struct {
		Config Config `json:"config"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &manifest.Config, nil
}
