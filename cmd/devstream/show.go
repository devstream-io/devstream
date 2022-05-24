package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/show/config"
	"github.com/devstream-io/devstream/internal/pkg/show/status"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var plugin string
var instanceID string
var statusAllFlag bool

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
  dtm show status --plugin=A-PLUGIN-NAME --id=A-PLUGIN-INSTANCE-ID
  dtm show status -p=A-PLUGIN-NAME -i=INSTANCE-ID
  dtm show status --all
  dtm show status -a`,
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

	showConfigCMD.Flags().StringVarP(&plugin, "plugin", "p", "", "specify name with the plugin")

	showStatusCMD.Flags().StringVarP(&plugin, "plugin", "p", "", "specify name with the plugin")
	showStatusCMD.Flags().StringVarP(&instanceID, "id", "i", "", "specify id with the plugin instance")
	showStatusCMD.Flags().BoolVarP(&statusAllFlag, "all", "a", false, "show all instances of all plugins status")
	showStatusCMD.Flags().StringVarP(&pluginDir, "plugin-dir", "d", pluginengine.DefaultPluginDir, "plugins directory")
	showStatusCMD.Flags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
}
