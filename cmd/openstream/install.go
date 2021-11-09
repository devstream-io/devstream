package main

import (
	"github.com/ironcore864/openstream/internal/pkg/config"
	"github.com/ironcore864/openstream/internal/pkg/plugin"
	"github.com/spf13/cobra"
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
		plugin.Install(&tool)
	}
}
