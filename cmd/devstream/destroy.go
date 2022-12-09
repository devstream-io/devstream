package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var isForceDestroy bool

var destroyCMD = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy DevOps tools deployment according to DevStream configuration file & state file",
	Long:  `Destroy DevOps tools deployment according to DevStream configuration file & state file.`,
	Run:   destroyCMDFunc,
}

func destroyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Destroy started.")
	if err := pluginengine.Destroy(configFilePath, continueDirectly, isForceDestroy); err != nil {
		log.Errorf("Destroy failed => %s.", err)
		os.Exit(1)
	}
	log.Success("Destroy finished.")
}

func init() {
	destroyCMD.Flags().BoolVarP(&isForceDestroy, "force", "", false, "force destroy by config")
	destroyCMD.Flags().StringVarP(&configFilePath, configFlagName, "f", "", "config file or directory")
	destroyCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", defaultPluginDir, "plugins directory")
	destroyCMD.Flags().BoolVarP(&continueDirectly, "yes", "y", false, "destroy directly without confirmation")

	completion.FlagFilenameCompletion(destroyCMD, configFlagName)
}
