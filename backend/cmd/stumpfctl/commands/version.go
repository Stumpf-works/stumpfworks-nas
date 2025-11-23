package commands

import (
	"github.com/Stumpf-works/stumpfworks-nas/pkg/cli"
	"github.com/Stumpf-works/stumpfworks-nas/pkg/client"
	"github.com/spf13/cobra"
)

// VersionCmd returns the version command
func VersionCmd(version, buildTime string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			cli.PrintHeader("StumpfWorks NAS")

			data := map[string]string{
				"CLI Version":    version,
				"Build Time":     buildTime,
				"Server Version": getServerVersion(),
			}

			cli.KeyValueTable(data)
		},
	}
}

func getServerVersion() string {
	apiClient := client.NewClient("http://localhost:8080")
	version, err := apiClient.GetVersion()
	if err != nil {
		return "unavailable (server not running)"
	}
	return version
}
