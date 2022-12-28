package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	configFilePath   string
	pluginDir        string
	continueDirectly bool
)

const (
	configFlagName    = "config-file"
	pluginDirFlagName = "plugin-dir"
	defaultPluginDir  = "~/.devstream/plugins"
)

func checkConfigFile() {
	if strings.TrimSpace(configFilePath) == "" {
		log.Errorf(`Config file is required. You could use "-f filename" or "-f directory" to specify it.`)
		os.Exit(1)
	}
}

func addFlagConfigFile(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&configFilePath, configFlagName, "f", "", "config file or directory")
	completion.FlagFilenameCompletion(cmd, configFlagName)
}

func addFlagPluginDir(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", defaultPluginDir, "plugins directory")
	completion.FlagDirnameCompletion(cmd, pluginDirFlagName)
}

func addFlagContinueDirectly(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&continueDirectly, "yes", "y", false, "continue directly without confirmation")
}
