package main

import (
	"fmt"
	"github.com/devstream-io/devstream/cmd/devstream/validator"

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
	Run: validator.WithValidators(showCMDFunc, validator.ArgsCountEqual(1), validateShowArgs),
}

func showCMDFunc(cmd *cobra.Command, args []string) {
	showInfo := show.Info(args[0])
	log.Debugf("The show info is: %s.", showInfo)
	if err := show.GenerateInfo(showInfo); err != nil {
		log.Fatal(err)
	}
}

func validateShowArgs(args []string) error {
	// arg is "config" or "status" here, maybe will have "output" in the future.
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
