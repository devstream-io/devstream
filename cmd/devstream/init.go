package main

import (
	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
)

var initCMD = &cobra.Command{
	Use:   "init",
	Short: "Download needed plugins according to the config file",
	Long:  `Download needed plugins according to the config file`,
	Run:   initCMDFunc,
}

func initCMDFunc(cmd *cobra.Command, args []string) {
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		log.Fatal("Failed to load the config file.")
	}

	log.Info("Initialize started.")
	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	log.Success("Initialize finished.")
}
