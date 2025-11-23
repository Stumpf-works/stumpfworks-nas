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

			apiClient := client.NewClient("http://localhost:8080")
			info, err := apiClient.GetSystemInfo()
			if err != nil {
				cli.PrintError("Failed to retrieve system information: %v", err)
				return err
			}

			// Display system information
			if hostname, ok := info["hostname"].(string); ok {
				fmt.Printf("Hostname: %s\n", hostname)
			}
			if uptime, ok := info["uptime"].(string); ok {
				fmt.Printf("Uptime: %s\n", uptime)
			}
			if cpuUsage, ok := info["cpu_usage"].(float64); ok {
				fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
			}
			if memUsage, ok := info["memory_usage"].(float64); ok {
				fmt.Printf("Memory Usage: %.2f%%\n", memUsage)
			}
			if diskUsage, ok := info["disk_usage"].(float64); ok {
				fmt.Printf("Disk Usage: %.2f%%\n", diskUsage)
			}

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
