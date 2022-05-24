package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var initCMD = &cobra.Command{
	Use:   "init",
	Short: "Download needed plugins according to the config file",
	Long:  `Download needed plugins according to the config file`,
	Run:   initCMDFunc,
}

func initCMDFunc(cmd *cobra.Command, args []string) {
	cfg, err := configloader.LoadConf(configFile)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	err = pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	log.Success("Initialize finished.")
}

func init() {
	initCMD.Flags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
	initCMD.Flags().StringVarP(&pluginDir, "plugin-dir", "d", pluginengine.DefaultPluginDir, "plugins directory")
}
