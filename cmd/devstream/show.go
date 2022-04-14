package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/show"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var plugin string
var instanceName string

var showCMD = &cobra.Command{
	Use:   "show [config | status]",
	Short: "Show is used to print some useful information",
	Long: `Show is used to print some useful information. 
Examples:
  dtm show config --plugin=A-PLUGIN-NAME
  dtm show status --plugin=A-PLUGIN-NAME --name=A-PLUGIN-INSTANCE-NAME
  dtm show status -p=A-PLUGIN-NAME -n=A-PLUGIN-INSTANCE-NAME`,
	Run: showCMDFunc,
}

func showCMDFunc(cmd *cobra.Command, args []string) {
	if err := validateShowArgs(args); err != nil {
		log.Fatal(err)
	}

	showInfo := show.Info(args[0])
	log.Debugf("The show info is: %s.", showInfo)
	if err := show.GenerateInfo(showInfo); err != nil {
		log.Fatal(err)
	}
}

func validateShowArgs(args []string) error {
	// arg is "config" or "status" here, maybe will have "output" in the future.
	if len(args) != 1 {
		return fmt.Errorf("got illegal args count (expect 1, got %d)", len(args))
	}
	showInfo := show.Info(args[0])
	if !show.IsValideInfo(showInfo) {
		return fmt.Errorf("invalide Show Info")
	}
	return nil
}

func init() {
	showCMD.PersistentFlags().StringVarP(&plugin, "plugin", "p", "", "specify name with the plugin")
	showCMD.PersistentFlags().StringVarP(&instanceName, "name", "n", "", "specify name with the plugin instance")
}
