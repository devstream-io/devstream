package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/list"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "This command lists all of the plugins",
	Long: `This command lists all of the plugins.
eg.
- dtm list plugins`,
	Run: listCMDFunc,
}

func listCMDFunc(cmd *cobra.Command, args []string) {
	if err := validateListCMDArgs(args); err != nil {
		log.Fatal(err)
	}

	listAction := list.Action(args[0])
	log.Debugf("The list action is: %s.", listAction)
	if err := list.ExecuteAction(listAction); err != nil {
		log.Fatal(err)
	}
}

func validateListCMDArgs(args []string) error {
	// "plugins"/ maybe it will be "core" in future.
	if len(args) != 1 {
		return fmt.Errorf("got illegal args count (expect 1, got %d). "+
			"See `help` command for more info", len(args))
	}
	listAction := list.Action(args[0])
	if !list.IsValideAction(listAction) {
		return fmt.Errorf("invalide Develop Action")
	}
	return nil
}

// TODO Use `--filter=someone` (can support regex) to filter plugins on feature,
// TODO Use `--group=somegroup` to filter the specified groups on feature
// func init() {
// 	developCMD.PersistentFlags().StringVarP(&filter, "filter", "f", "", "filter plugins; support regex ")
// 	developCMD.PersistentFlags().StringVarP(&group, "group", "g", "", "filter the specified groups; support regex")
// }
