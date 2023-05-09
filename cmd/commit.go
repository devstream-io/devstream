package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/log"
	"github.com/devstream-io/devstream/internal/pkg/commit"
	"github.com/devstream-io/devstream/internal/response"
)

// commit message got from the command line by -m flag
var message string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit is used to execute git commit operations",
	Long: `commit is used to execute git commit operations

e.g.

1. dtm commit -m "commit message"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if message == "" {
			errStr := "message is required"
			log.Error(errStr)
			r := response.New(response.StatusError, response.MessageError, errStr)
			r.Print(OutputFormat)
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
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
}
