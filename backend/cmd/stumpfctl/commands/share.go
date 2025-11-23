package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// ShareCmd returns the share management command
func ShareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Manage shares",
		Long:  "Create, delete, and list Samba shares",
	}

	cmd.AddCommand(shareListCmd())

	return cmd
}

func shareListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all shares",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.PrintHeader("StumpfWorks NAS Shares")

			apiClient := client.NewClient("http://localhost:8080")
			shares, err := apiClient.GetShares()
			if err != nil {
				cli.PrintError("Failed to retrieve shares: %v", err)
				return err
			}

			if len(shares) == 0 {
				fmt.Println("No shares configured")
				return nil
			}

			// Display shares
			for _, share := range shares {
				if name, ok := share["name"].(string); ok {
					fmt.Printf("\nShare: %s\n", name)
				}
				if path, ok := share["path"].(string); ok {
					fmt.Printf("  Path: %s\n", path)
				}
				if shareType, ok := share["type"].(string); ok {
					fmt.Printf("  Type: %s\n", shareType)
				}
				if enabled, ok := share["enabled"].(bool); ok {
					fmt.Printf("  Enabled: %v\n", enabled)
				}
			}

			return nil
		},
	}
}
