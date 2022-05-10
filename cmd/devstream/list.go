package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	pluginFilter string
)

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "This command lists all of the plugins",
	Long: `This command lists all of the plugins.
Examples:
  dtm list plugins`,
	Run: listCMDFunc,
}

func listCMDFunc(cmd *cobra.Command, args []string) {
	if err := validateListCMDArgs(args); err != nil {
		log.Fatal(err)
	}

	list.List(pluginFilter)
}

func validateListCMDArgs(args []string) error {
	// only support "plugins" now
	if len(args) != 1 {
		return fmt.Errorf("got illegal args count (expect 1, got %d)", len(args))
	}

	if args[0] != "plugins" {
		return fmt.Errorf("arg should be \"plugins\" only")
	}
	return nil
}

// TODO Use `--group=somegroup` to filter the specified groups on feature
func init() {
	listCMD.PersistentFlags().StringVarP(&pluginFilter, "filter", "r", "", "filter plugin by regex")
}
