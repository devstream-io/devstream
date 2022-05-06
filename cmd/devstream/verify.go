package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/configloader"

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

	gConfig, err := configloader.LoadGeneralConf(configFile)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}
	log.Debugf("config file content is %s.", gConfig)

	if pluginengine.Verify(gConfig.ToolFile, gConfig.VarFile) {
		log.Success("Verify succeeded.")
	} else {
		log.Info("Verify finished.")
	}
}
