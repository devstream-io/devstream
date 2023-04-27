package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/log"
	"github.com/devstream-io/devstream/internal/pkg/commit"
	"github.com/devstream-io/devstream/internal/response"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit is used to execute git commit operations",
	Long: `commit is used to execute git commit operations

e.g.

1. dtm commit -m "commit message"
`,
	Run: func(cmd *cobra.Command, args []string) {
		message := viper.GetString("message")
		if message == "" {
			log.Error("message is required")
			os.Exit(1)
		}
		err := commit.Commit(message)
		if err != nil {
			log.Errorf("commit error: %v", err)
			r := response.New(response.StatusError, response.MessageError, err.Error())
			r.Print(OutputFormat)
		} else {
			r := response.New(response.StatusOK, response.MessageOK, "")
			r.Print(OutputFormat)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringP("message", "m", "", "commit message")
}
