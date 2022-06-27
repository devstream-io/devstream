package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var isForceDelete bool

var deleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Delete DevOps tools according to DevStream configuration file",
	Long: `Delete DevOps tools according to DevStream configuration file.
DevStream will delete everything defined in the config file, regardless of the state.`,
	Run: deleteCMDFunc,
}

func deleteCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Delete started.")
	if err := pluginengine.Remove(configFile, continueDirectly, isForceDelete); err != nil {
		log.Errorf("Delete error: %s.", err)
		os.Exit(1)
	}

	log.Success("Delete finished.")
}

func init() {
	deleteCMD.Flags().BoolVarP(&isForceDelete, "force", "", false, "force delete by config")
	deleteCMD.Flags().StringVarP(&configFile, configFlagName, "f", "config.yaml", "config file")
	deleteCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", pluginengine.DefaultPluginDir, "plugins directory")
	deleteCMD.Flags().BoolVarP(&continueDirectly, "yes", "y", false, "delete directly without confirmation")

	completion.FlagFilenameCompletion(deleteCMD, configFlagName)
	completion.FlagDirnameCompletion(deleteCMD, pluginDirFlagName)
}
