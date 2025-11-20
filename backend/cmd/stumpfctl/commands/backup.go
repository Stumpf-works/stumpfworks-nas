package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// BackupCmd returns the backup management command
func BackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage backups",
		Long:  "Create, list, and restore database backups",
	}

	cmd.AddCommand(backupListCmd())
	cmd.AddCommand(backupCreateCmd())

	return cmd
}

func backupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient := client.NewClient("http://localhost:8080")

			backups, err := apiClient.GetBackups()
			if err != nil {
				cli.PrintError("Failed to retrieve backups: %v", err)
				return err
			}

			cli.PrintHeader("StumpfWorks NAS Backups")

			headers := []string{"Filename", "Size", "Created"}
			rows := [][]string{}

			for _, backup := range backups {
				filename := fmt.Sprintf("%v", backup["filename"])
				size := fmt.Sprintf("%v", backup["size"])
				created := fmt.Sprintf("%v", backup["created"])

				rows = append(rows, []string{filename, size, created})
			}

			cli.Table(headers, rows)
			fmt.Printf("\nTotal: %d backups\n", len(backups))

			return nil
		},
	}
}

func backupCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a new backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintInfo("Creating backup...")

			apiClient := client.NewClient("http://localhost:8080")
			if err := apiClient.CreateBackup(); err != nil {
				cli.PrintError("Failed to create backup: %v", err)
				return err
			}

			cli.PrintSuccess("Backup created successfully")
			return nil
		},
	}
}
