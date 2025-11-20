package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// HealthCmd returns the health command
func HealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check system health",
		Long:  "Perform a comprehensive health check of StumpfWorks NAS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHealth()
		},
	}
}

func checkHealth() error {
	cli.PrintHeader("StumpfWorks NAS Health Check")

	// Create API client
	apiClient := client.NewClient("http://localhost:8080")

	// Check API health
	cli.PrintInfo("Checking API health...")
	health, err := apiClient.Health()
	if err != nil {
		cli.PrintError("API is not responding: %v", err)
		return err
	}

	cli.PrintSuccess("API is healthy")

	// Display health data
	fmt.Println()
	fmt.Println("Health Report:")
	for key, value := range health {
		fmt.Printf("  %s: %v\n", key, value)
	}

	return nil
}
