package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// SystemCmd returns the system information command
func SystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System information and maintenance",
		Long:  "View system information, metrics, and perform maintenance tasks",
	}

	cmd.AddCommand(systemInfoCmd())
	cmd.AddCommand(systemMetricsCmd())

	return cmd
}

func systemInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show system information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintHeader("StumpfWorks NAS System Information")
			fmt.Println("System info display not yet implemented")
			return nil
		},
	}
}

func systemMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metrics",
		Short: "Show system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient := client.NewClient("http://localhost:8080")

			metrics, err := apiClient.GetMetrics()
			if err != nil {
				cli.PrintError("Failed to retrieve metrics: %v", err)
				return err
			}

			cli.PrintHeader("StumpfWorks NAS System Metrics")

			for key, value := range metrics {
				fmt.Printf("%s: %v\n", key, value)
			}

			return nil
		},
	}
}
