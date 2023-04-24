package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/git"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "git is used to execute git operations",
	Long: `git is used to execute git operations
1. dtm git 'git add .'
2. dtm git 'git commit -s -m "commit message"'
`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCmdStr := strings.Join(args, " ")
		git.Execute(gitCmdStr)
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
