package main

import (
	"log"

	"github.com/spf13/cobra"

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

	log.Println("Initialize started.")
	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	log.Println("Initialize finished.")
}
