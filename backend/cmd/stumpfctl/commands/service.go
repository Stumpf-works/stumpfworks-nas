package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/spf13/cobra"
)

const serviceName = "stumpfworks-nas.service"

// ServiceCmd returns the service management command
func ServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage StumpfWorks NAS service",
		Long:  "Start, stop, restart, and check status of the StumpfWorks NAS service",
	}

	cmd.AddCommand(startCmd())
	cmd.AddCommand(stopCmd())
	cmd.AddCommand(restartCmd())
	cmd.AddCommand(statusCmd())
	cmd.AddCommand(enableCmd())
	cmd.AddCommand(disableCmd())
	cmd.AddCommand(reloadCmd())

	return cmd
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the StumpfWorks NAS service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Starting StumpfWorks NAS service...")
			return runSystemctl("start")
		},
	}
}

func stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the StumpfWorks NAS service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cli.ConfirmPrompt("Are you sure you want to stop the service?") {
				cli.PrintInfo("Cancelled")
				return nil
			}
			cli.PrintInfo("Stopping StumpfWorks NAS service...")
			return runSystemctl("stop")
		},
	}
}

func restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the StumpfWorks NAS service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Restarting StumpfWorks NAS service...")
			return runSystemctl("restart")
		},
	}
}

func reloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reload",
		Short: "Reload the StumpfWorks NAS configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Reloading StumpfWorks NAS configuration...")
			return runSystemctl("reload")
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showDetailedStatus()
		},
	}
}

func enableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable auto-start on boot",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Enabling auto-start...")
			return runSystemctl("enable")
		},
	}
}

func disableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable auto-start on boot",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Disabling auto-start...")
			return runSystemctl("disable")
		},
	}
}

// runSystemctl executes a systemctl command
func runSystemctl(action string) error {
	cmd := exec.Command("systemctl", action, serviceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		cli.PrintError("Failed to %s service: %s", action, string(output))
		return err
	}

	cli.PrintSuccess("Service %s successful", action)
	return nil
}

// showDetailedStatus shows a detailed status of the service
func showDetailedStatus() error {
	cli.PrintHeader("StumpfWorks NAS Status")

	// Check if service is running
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, _ := cmd.Output()
	isActive := strings.TrimSpace(string(output)) == "active"

	// Get service status
	cmd = exec.Command("systemctl", "status", serviceName, "--no-pager", "-l")
	output, _ = cmd.CombinedOutput()
	statusOutput := string(output)

	// Extract uptime
	var uptimeStr string
	if isActive {
		cmd = exec.Command("systemctl", "show", serviceName, "--property=ActiveEnterTimestamp")
		output, _ = cmd.Output()
		uptimeStr = strings.TrimSpace(strings.TrimPrefix(string(output), "ActiveEnterTimestamp="))
	}

	// Build status display
	status := "● Running (healthy)"
	if !isActive {
		status = cli.Error("✗ Stopped")
	}

	data := map[string]string{
		"Service":  status,
		"Uptime":   uptimeStr,
		"Version":  "v0.1.0",
	}

	cli.KeyValueTable(data)

	fmt.Println()
	fmt.Println("Systemd Status:")
	fmt.Println(statusOutput)

	return nil
}
