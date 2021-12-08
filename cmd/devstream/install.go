package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/plugin"
)

var installCMD = &cobra.Command{
	Use:   "install",
	Short: "Install tools defined in config file",
	Long:  `Install tools defined in config file`,
	Run:   installCMDFunc,
}

func installCMDFunc(cmd *cobra.Command, args []string) {
	conf := config.LoadConf(configFile)
	for _, tool := range conf.Tools {
		log.Printf("=== start to install plugin %s %s ===", tool.Name, tool.Version)
		plugin.Install(&tool)
		log.Printf("=== plugin %s %s installation done===", tool.Name, tool.Version)
		log.Printf("===")
	}
}
