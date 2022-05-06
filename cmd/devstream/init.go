package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
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
	gConfig, err := configloader.LoadGeneralConf(configFile)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}
	log.Debugf("config file content is %s.", gConfig)

	cfg := configloader.LoadConf(gConfig.ToolFile, gConfig.VarFile)
	if cfg == nil {
		log.Fatal("Failed to load the config file.")
	}

	err = pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	log.Success("Initialize finished.")
}
