package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/show/config"
	"github.com/devstream-io/devstream/internal/pkg/show/status"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var plugin string
var instanceName string

var showCMD = &cobra.Command{
	Use:   "show",
	Short: "Show is used to print some useful information",
}

var showConfigCMD = &cobra.Command{
	Use:   "config",
	Short: "Show configuration information",
	Long: `Show config is used for showing plugins' template configuration information.
Examples:
  dtm show config --plugin=A-PLUGIN-NAME`,
	Run: showConfigCMDFunc,
}

var showStatusCMD = &cobra.Command{
	Use:   "status",
	Short: "Show status information",
	Long: `Show status is used for showing plugins' status information.
Examples:
  dtm show status --plugin=A-PLUGIN-NAME --name=A-PLUGIN-INSTANCE-NAME
  dtm show status -p=A-PLUGIN-NAME -n=A-PLUGIN-INSTANCE-NAME`,
	Run: showStatusCMDFunc,
}

func showConfigCMDFunc(cmd *cobra.Command, args []string) {
	log.Debug("Show configuration information.")
	if err := config.Show(); err != nil {
		log.Fatal(err)
	}
}

func showStatusCMDFunc(cmd *cobra.Command, args []string) {
	log.Debug("Show status information.")
	if err := status.Show(configFile); err != nil {
		log.Fatal(err)
	}
}

func init() {
	showCMD.AddCommand(showConfigCMD)
	showCMD.AddCommand(showStatusCMD)

	showConfigCMD.PersistentFlags().StringVarP(&plugin, "plugin", "p", "", "specify name with the plugin")
	showStatusCMD.PersistentFlags().StringVarP(&plugin, "plugin", "p", "", "specify name with the plugin")
	showStatusCMD.PersistentFlags().StringVarP(&instanceName, "name", "n", "", "specify name with the plugin instance")
}
