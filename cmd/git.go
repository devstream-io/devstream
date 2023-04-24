package cmd

import (
	"github.com/devstream-io/devstream/internal/pkg/git"
	"github.com/spf13/cobra"
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
