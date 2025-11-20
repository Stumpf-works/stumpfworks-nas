package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
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
			fmt.Println("Share listing not yet implemented")
			return nil
		},
	}
}
