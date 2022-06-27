package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/completion"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var verifyCMD = &cobra.Command{
	Use:   "verify",
	Short: "Verify DevOps tools according to DevStream config file and state",
	Long:  `Verify DevOps tools according to DevStream config file and state.`,
	Run:   verifyCMDFunc,
}

func verifyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Verify started.")
	if pluginengine.Verify(configFile) {
		log.Success("Verify succeeded.")
	} else {
		log.Info("Verify finished.")
	}
}

func init() {
	verifyCMD.Flags().StringVarP(&configFile, configFlagName, "f", "config.yaml", "config file")
	verifyCMD.Flags().StringVarP(&pluginDir, pluginDirFlagName, "d", pluginengine.DefaultPluginDir, "plugins directory")

	completion.FlagFilenameCompletion(verifyCMD, configFlagName)
	completion.FlagDirnameCompletion(verifyCMD, pluginDirFlagName)
}
