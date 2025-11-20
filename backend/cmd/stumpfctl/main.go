package main

import (
	"fmt"
	"os"

	"github.com/Stumpf-works/stumpfworks-nas/cmd/stumpfctl/commands"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "stumpfctl",
		Short: "StumpfWorks NAS Management CLI",
		Long: `stumpfctl is the command-line management tool for StumpfWorks NAS.
It provides a user-friendly interface to manage services, users, backups,
configuration, shares, and monitor system health.`,
		Version: fmt.Sprintf("%s (built %s)", Version, BuildTime),
	}

	// Add all subcommands
	rootCmd.AddCommand(commands.ServiceCmd())
	rootCmd.AddCommand(commands.LogsCmd())
	rootCmd.AddCommand(commands.UserCmd())
	rootCmd.AddCommand(commands.BackupCmd())
	rootCmd.AddCommand(commands.ConfigCmd())
	rootCmd.AddCommand(commands.ShareCmd())
	rootCmd.AddCommand(commands.HealthCmd())
	rootCmd.AddCommand(commands.SystemCmd())
	rootCmd.AddCommand(commands.VersionCmd(Version, BuildTime))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
