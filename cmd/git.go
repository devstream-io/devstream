package cmd

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/git"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "git is used to execute git operations",
	Long:  `git is used to execute git operations`,
	Run: func(cmd *cobra.Command, args []string) {
		git.Run()
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
