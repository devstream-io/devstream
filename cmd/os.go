package cmd

import (
	"github.com/devstream-io/devstream/internal/pkg/osx"
	"github.com/spf13/cobra"
)

// osCmd represents the os command
var osCmd = &cobra.Command{
	Use:   "os",
	Short: "os is used to execute os operations",
	Long:  `os is used to execute os operations`,
	Run: func(cmd *cobra.Command, args []string) {
		osx.Run()
	},
}

func init() {
	rootCmd.AddCommand(osCmd)
}
