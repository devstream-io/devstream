package main

import (
	"log"

	"github.com/spf13/viper"

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
	var tools = make([]configloader.Tool, 0)
	if err := viper.UnmarshalKey("tools", &tools); err != nil {
		log.Fatal(err)
	}
	var cfg = &configloader.Config{
		Tools: tools,
	}

	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	log.Println("=== initialize finished ===")
}
