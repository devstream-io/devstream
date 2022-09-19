package main

import (
	"fmt"
	"runtime"
	"strings"

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

	switch downloadOnly {
	// download plugins according to the config file
	case false:
		tools, realPluginDir, err = GetPluginsAndPluginDirFromConfig()
	// download plugins from flags
	case true:
		tools, realPluginDir, err = GetPluginsAndPluginDirFromFlags()
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

func GetPluginsAndPluginDirFromConfig() (tools []configmanager.Tool, pluginName string, err error) {
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		return nil, "", err
	}

	// combine plugin dir from config file and flag
	if err := file.SetPluginDir(cfg.PluginDir); err != nil {
		return nil, "", err
	}

	return cfg.Tools, viper.GetString(pluginDirFlagName), nil
}

func GetPluginsAndPluginDirFromFlags() (tools []configmanager.Tool, pluginName string, err error) {
	// 1. get plugins from flags
	var pluginsName []string
	switch downloadAll {
	// download all plugins
	case true:
		pluginsName = list.PluginsNameSlice()
	// download specific plugins
	case false:
		log.Info(pluginsToDownload)
		for _, pluginName := range pluginsToDownload {
			if p := strings.ToLower(strings.TrimSpace(pluginName)); p != "" {
				pluginsName = append(pluginsName, p)
			}
		}
		// check if plugins to download are supported by dtm
		for _, plugin := range pluginsToDownload {
			if _, ok := list.PluginNamesMap()[plugin]; !ok {
				return nil, "", fmt.Errorf("plugin %s is not supported by dtm", plugin)
			}
		}
	}

	if len(pluginsName) == 0 {
		log.Errorf("Please use --plugins to specify plugins to download or use --all to download all plugins.")
	}
	log.Debugf("plugins to download: %v", pluginsName)

	if initOS == "" || initArch == "" {
		return nil, "", fmt.Errorf("once you use the --all flag, you must specify the --os and --arch flags")
	}

	// build the plugin list
	for _, pluginName := range pluginsName {
		tools = append(tools, configmanager.Tool{Name: pluginName})
	}

	// 2. handle plugin dir
	pluginDir = viper.GetString(pluginDirFlagName)
	pluginDir, err = file.HandlePathWithHome(pluginDir)
	if err != nil {
		return nil, "", err
	}

	return tools, pluginDir, nil
}

func init() {
	// flags for init from config file
	initCMD.Flags().StringVarP(&configFile, configFlagName, "f", "config.yaml", "config file")
	initCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", "", "plugins directory")

	// downloading specific plugins from flags
	initCMD.Flags().BoolVar(&downloadOnly, "download-only", false, "download plugins only")
	initCMD.Flags().StringSliceVarP(&pluginsToDownload, "plugins", "p", []string{}, "plugins to download")
	initCMD.Flags().BoolVarP(&downloadAll, "all", "a", false, "download all plugins")
	initCMD.Flags().StringVar(&initOS, "os", runtime.GOOS, "download plugins for specific os")
	initCMD.Flags().StringVar(&initArch, "arch", runtime.GOARCH, "download plugins for specific arch")

	completion.FlagFilenameCompletion(initCMD, configFlagName)
	completion.FlagDirnameCompletion(initCMD, pluginDirFlagName)
}
