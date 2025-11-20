package commands

import (
	"fmt"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// UserCmd returns the user management command
func UserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Long:  "Create, delete, and list StumpfWorks NAS users",
	}

	cmd.AddCommand(userListCmd())
	cmd.AddCommand(userAddCmd())
	cmd.AddCommand(userDeleteCmd())

	return cmd
}

func userListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all users",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient := client.NewClient("http://localhost:8080")

			users, err := apiClient.GetUsers()
			if err != nil {
				cli.PrintError("Failed to retrieve users: %v", err)
				return err
			}

			cli.PrintHeader("StumpfWorks NAS Users")

			headers := []string{"Username", "Role", "Status", "Last Login"}
			rows := [][]string{}

			for _, user := range users {
				username := fmt.Sprintf("%v", user["username"])
				role := fmt.Sprintf("%v", user["role"])
				status := "Active"
				lastLogin := "Never"

				if ll, ok := user["last_login"]; ok && ll != nil {
					lastLogin = fmt.Sprintf("%v", ll)
				}

				rows = append(rows, []string{username, role, status, lastLogin})
			}

			cli.Table(headers, rows)
			fmt.Printf("\nTotal: %d users\n", len(users))

			return nil
		},
	}
}

func userAddCmd() *cobra.Command {
	var (
		admin bool
	)

	cmd := &cobra.Command{
		Use:   "add <username>",
		Short: "Add a new user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			// Prompt for password
			password, err := cli.PasswordPrompt("Enter password")
			if err != nil {
				return err
			}

			passwordConfirm, err := cli.PasswordPrompt("Confirm password")
			if err != nil {
				return err
			}

			if password != passwordConfirm {
				cli.PrintError("Passwords do not match")
				return fmt.Errorf("passwords do not match")
			}

			role := "user"
			if admin {
				role = "admin"
			}

			apiClient := client.NewClient("http://localhost:8080")
			if err := apiClient.CreateUser(username, password, role); err != nil {
				cli.PrintError("Failed to create user: %v", err)
				return err
			}

			cli.PrintSuccess("User '%s' created successfully", username)
			return nil
		},
	}

	cmd.Flags().BoolVar(&admin, "admin", false, "Create as admin user")

	return cmd
}

func userDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <username>",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]

			if !cli.ConfirmPrompt(fmt.Sprintf("Are you sure you want to delete user '%s'?", username)) {
				cli.PrintInfo("Cancelled")
				return nil
			}

			apiClient := client.NewClient("http://localhost:8080")
			if err := apiClient.DeleteUser(username); err != nil {
				cli.PrintError("Failed to delete user: %v", err)
				return err
			}

			cli.PrintSuccess("User '%s' deleted successfully", username)
			return nil
		},
	}
}
