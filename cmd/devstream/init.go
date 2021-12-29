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
	conf := configloader.LoadConf(configFile)
	err := pluginmanager.DownloadPlugins(conf)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	log.Println("=== initialize finished ===")
}
