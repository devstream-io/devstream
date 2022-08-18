package main

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/file"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var initCMD = &cobra.Command{
	Use:   "init",
	Short: "Download needed plugins according to the config file",
	Long:  `Download needed plugins according to the config file`,
	Run:   initCMDFunc,
}

func initCMDFunc(_ *cobra.Command, _ []string) {
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	file.SetPluginDir(cfg.PluginDir)

	if version.Dev {
		log.Errorf("Dev version plugins can't be downloaded from the remote plugin repo; please run `make build-plugin.PLUGIN_NAME` to build them locally.")
		return
	}

	if err = pluginmanager.DownloadPlugins(cfg); err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	fmt.Println()
	log.Success("Initialize finished.")
}

func init() {
	initCMD.Flags().StringVarP(&configFile, configFlagName, "f", "config.yaml", "config file")
	initCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", "", "plugins directory")

	completion.FlagFilenameCompletion(initCMD, configFlagName)
	completion.FlagDirnameCompletion(initCMD, pluginDirFlagName)
}
