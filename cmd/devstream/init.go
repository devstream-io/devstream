package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var initCMD = &cobra.Command{
	Use:   "init",
	Short: "Download needed plugins according to the config file",
	Long:  `Download needed plugins according to the config file`,
	Run:   initCMDFunc,
}

var (
	downloadAll bool   // download all plugins which dtm supports
	initOS      string // download plugins for specific os
	initArch    string // download plugins for specific arch
)

func initCMDFunc(_ *cobra.Command, _ []string) {
	if version.Dev {
		log.Errorf("Dev version plugins can't be downloaded from the remote plugin repo; please run `make build-plugin.PLUGIN_NAME` to build them locally.")
		return
	}

	var tools []configmanager.Tool

	switch downloadAll {
	case false:
		// download plugins according to the config file
		cfg, err := configmanager.NewManager(configFile).LoadConfig()
		if err != nil {
			log.Errorf("Error: %s.", err)
			return
		}

		// combine plugin dir from config file and flag
		if err := file.SetPluginDir(cfg.PluginDir); err != nil {
			log.Errorf("Error: %s.", err)
		}

		tools = cfg.Tools
		pluginDir = viper.GetString(pluginDirFlagName)
	case true:
		// download all plugins
		if initOS == "" || initArch == "" {
			log.Errorf("Once you use the --all flag, you must specify the --os and --arch flags.")
			return
		}
		// build the plugin list
		pluginsName := list.PluginsNameSlice()
		for _, pluginName := range pluginsName {
			tools = append(tools, configmanager.Tool{Name: pluginName})
		}

		var err error
		pluginDir = viper.GetString(pluginDirFlagName)
		pluginDir, err = file.HandlePathWithHome(pluginDir)
		if err != nil {
			log.Errorf("Error: %s.", err)
		}
	}

	if err := pluginmanager.DownloadPlugins(tools, pluginDir, initOS, initArch); err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	fmt.Println()
	log.Success("Initialize finished.")
}

func init() {
	initCMD.Flags().StringVarP(&configFile, configFlagName, "f", "config.yaml", "config file")
	initCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", "", "plugins directory")
	initCMD.Flags().BoolVarP(&downloadAll, "all", "a", false, "download all plugins which dtm supports")
	initCMD.Flags().StringVarP(&initOS, "os", "", runtime.GOOS, "download plugins for specific os")
	initCMD.Flags().StringVarP(&initArch, "arch", "", runtime.GOARCH, "download plugins for specific arch")

	completion.FlagFilenameCompletion(initCMD, configFlagName)
	completion.FlagDirnameCompletion(initCMD, pluginDirFlagName)
}
