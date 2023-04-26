package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/github"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "github is used to execute github operations",
	Long: `github is used to execute github operations
	参考 gh`,
	Run: func(cmd *cobra.Command, args []string) {
		github.Run()
	},
}

func init() {
	rootCmd.AddCommand(githubCmd)
}
