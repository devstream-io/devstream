package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
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

var (
	downloadOnly      bool     // download plugins only, from command line flags
	downloadAll       bool     // download all plugins
	pluginsToDownload []string // download specific plugins
	initOS            string   // download plugins for specific os
	initArch          string   // download plugins for specific arch
)

func initCMDFunc(_ *cobra.Command, _ []string) {
	if version.Dev {
		log.Fatalf("Dev version plugins can't be downloaded from the remote plugin repo; please run `make build-plugin.PLUGIN_NAME` to build them locally.")
	}

	var (
		realPluginDir string
		tools         []configmanager.Tool
		err           error
	)

	if downloadOnly {
		// download plugins from flags
		tools, realPluginDir, err = GetPluginsAndPluginDirFromFlags()
	} else {
		// download plugins according to the config file
		tools, realPluginDir, err = GetPluginsAndPluginDirFromConfig()
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := pluginmanager.DownloadPlugins(tools, realPluginDir, initOS, initArch); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	log.Success("Initialize finished.")
}

func GetPluginsAndPluginDirFromConfig() (tools []configmanager.Tool, pluginDir string, err error) {
	cfg, err := configmanager.NewManager(configFilePath).LoadConfig()
	if err != nil {
		return nil, "", err
	}

	// combine plugin dir from config file and flag
	if err = pluginengine.SetPluginDir(cfg.PluginDir); err != nil {
		return nil, "", err
	}

	return cfg.Tools, viper.GetString(pluginDirFlagName), nil
}

func GetPluginsAndPluginDirFromFlags() (tools []configmanager.Tool, pluginDir string, err error) {
	// 1. get plugins from flags
	var pluginsName []string
	if downloadAll {
		// download all plugins
		pluginsName = list.PluginsNameSlice()
	} else {
		// download specific plugins
		for _, pluginName := range pluginsToDownload {
			if p := strings.ToLower(strings.TrimSpace(pluginName)); p != "" {
				pluginsName = append(pluginsName, p)
			}
		}
		// check if plugins to download are supported by dtm
		for _, plugin := range pluginsName {
			if _, ok := list.PluginNamesMap()[plugin]; !ok {
				return nil, "", fmt.Errorf("Plugin %s is not supported by dtm", plugin)
			}
		}
	}

	if len(pluginsName) == 0 {
		return nil, "", errors.New("Please use --plugins to specify plugins to download or use --all to download all plugins.")
	}
	log.Debugf("plugins to download: %v", pluginsName)

	if initOS == "" || initArch == "" {
		return nil, "", fmt.Errorf("Once you use the --all flag, you must specify the --os and --arch flags")
	}

	log.Infof("Plugins to download: %v", pluginsName)

	// build the plugin list
	for _, pluginName := range pluginsName {
		tools = append(tools, configmanager.Tool{Name: pluginName})
	}

	// 2. handle plugin dir
	if err = pluginengine.SetPluginDir(""); err != nil {
		return nil, "", err
	}

	return tools, viper.GetString(pluginDirFlagName), nil
}

func init() {
	// flags for init from config file
	initCMD.Flags().StringVarP(&configFilePath, configFlagName, "f", "config.yaml", "config file")
	initCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", "", "plugins directory")

	// downloading specific plugins from flags
	initCMD.Flags().BoolVar(&downloadOnly, "download-only", false, "download plugins only")
	initCMD.Flags().StringSliceVarP(&pluginsToDownload, "plugins", "p", []string{}, "the plugins to be downloaded")
	initCMD.Flags().BoolVarP(&downloadAll, "all", "a", false, "download all plugins")
	initCMD.Flags().StringVar(&initOS, "os", runtime.GOOS, "download plugins for specific os")
	initCMD.Flags().StringVar(&initArch, "arch", runtime.GOARCH, "download plugins for specific arch")

	completion.FlagFilenameCompletion(initCMD, configFlagName)
	completion.FlagDirnameCompletion(initCMD, pluginDirFlagName)
}
