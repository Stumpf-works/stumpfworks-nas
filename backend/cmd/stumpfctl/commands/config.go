package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/spf13/cobra"
)

// ConfigCmd returns the configuration management command
func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "View and modify StumpfWorks NAS configuration",
	}

	cmd.AddCommand(configShowCmd())
	cmd.AddCommand(configEditCmd())

	return cmd
}

func configShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintHeader("StumpfWorks NAS Configuration")
			fmt.Println("Configuration display not yet implemented")
			return nil
		},
	}
}

func configEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Opening configuration editor...")
			fmt.Println("Configuration editor not yet implemented")
			return nil
		},
	}
}
