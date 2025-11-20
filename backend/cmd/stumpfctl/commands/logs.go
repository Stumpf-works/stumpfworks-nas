package commands

import (
	"fmt"
	"os/exec"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/spf13/cobra"
)

// LogsCmd returns the logs command
func LogsCmd() *cobra.Command {
	var (
		follow bool
		lines  int
		since  string
	)

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View StumpfWorks NAS logs",
		Long:  "Display logs from the StumpfWorks NAS service using journalctl",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showLogs(follow, lines, since)
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output (like tail -f)")
	cmd.Flags().IntVarP(&lines, "lines", "n", 50, "Number of lines to show")
	cmd.Flags().StringVar(&since, "since", "", "Show logs since (e.g., '1h ago', '2025-01-01')")

	return cmd
}

func showLogs(follow bool, lines int, since string) error {
	args := []string{"-u", serviceName, "--no-pager"}

	if follow {
		args = append(args, "-f")
	} else {
		args = append(args, fmt.Sprintf("-n%d", lines))
	}

	if since != "" {
		args = append(args, "--since", since)
	}

	cli.PrintInfo("Showing logs for StumpfWorks NAS...")
	fmt.Println()

	cmd := exec.Command("journalctl", args...)
	cmd.Stdout = cmd.Stdout
	cmd.Stderr = cmd.Stderr

	// If following, run interactively
	if follow {
		return cmd.Run()
	}

	// Otherwise, just show output
	output, err := cmd.CombinedOutput()
	if err != nil {
		cli.PrintError("Failed to retrieve logs: %v", err)
		return err
	}

	fmt.Println(string(output))
	return nil
}
