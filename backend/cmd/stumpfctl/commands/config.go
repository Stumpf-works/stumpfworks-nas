package commands

import (
	"fmt"
	"os"
	"os/exec"

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

			configPath := "/etc/stumpfworks/config.yaml"
			data, err := os.ReadFile(configPath)
			if err != nil {
				cli.PrintError("Failed to read configuration file: %v", err)
				cli.PrintInfo("Configuration file location: %s", configPath)
				return err
			}

			fmt.Println(string(data))
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

			configPath := "/etc/stumpfworks/config.yaml"

			// Check if file exists
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				cli.PrintError("Configuration file not found: %s", configPath)
				return err
			}

			// Use EDITOR environment variable or default to vi
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}

			// Open editor
			editorCmd := exec.Command(editor, configPath)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			if err := editorCmd.Run(); err != nil {
				cli.PrintError("Failed to open editor: %v", err)
				return err
			}

			cli.PrintSuccess("Configuration file updated")
			cli.PrintInfo("Remember to restart the service for changes to take effect:")
			cli.PrintInfo("  sudo systemctl restart stumpfworks-server")
			return nil
		},
	}
}
