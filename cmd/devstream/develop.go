package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/develop"
	"github.com/merico-dev/stream/pkg/util/log"
)

var name string

var developCMD = &cobra.Command{
	Use:   "develop",
	Short: "Develop is used for develop a new plugin",
	Long: `Develop is used for develop a new plugin.
eg.
- dtm develop create-plugin --name=YOUR-PLUGIN-NAME`,
	Run: developCMDFunc,
}

func developCMDFunc(cmd *cobra.Command, args []string) {
	if err := validateArgs(args); err != nil {
		log.Fatal(err)
	}

	developAction := develop.Action(args[0])
	log.Debugf("The develop action is: %s", developAction)
	if err := develop.BranchAction(developAction); err != nil {
		log.Fatal(err)
	}
}

func validateArgs(args []string) error {
	// "create-plugin"/ maybe it will be "delete-plugin"/"rename-plugin" in future.
	if len(args) != 1 {
		return fmt.Errorf("got illegal args count (expect 1, got %d). "+
			"See `help` command for more info", len(args))
	}
	developAction := develop.Action(args[0])
	if !develop.IsValideAction(developAction) {
		return fmt.Errorf("invalide Develop Action")
	}
	return nil
}

func init() {
	developCMD.PersistentFlags().StringVarP(&name, "name", "n", "", "specify name with the new plugin")
}
