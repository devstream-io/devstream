package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var applyCMD = &cobra.Command{
	Use:   "apply",
	Short: "Create or update DevOps tools according to DevStream configuration file",
	Long: `Create or update DevOps tools according to DevStream configuration file.
DevStream will generate and execute a new plan based on the config file and the state file by default.`,
	Run:        applyCMDFunc,
	SuggestFor: []string{"install"},
}

func applyCMDFunc(cmd *cobra.Command, args []string) {
	checkConfigFile()
	log.Info("Apply started.")
	if err := pluginengine.Apply(configFilePath, continueDirectly); err != nil {
		log.Errorf("Apply failed => %s.", err)
		if strings.Contains(err.Error(), "config not valid") {
			log.Info("It seems your config file is not valid. Please check the official documentation https://docs.devstream.io, or use the \"dtm show config\" command to get an example.")
		}
		os.Exit(1)
	}
	log.Success("Apply finished.")
}

func init() {
	addFlagConfigFile(applyCMD)
	addFlagPluginDir(applyCMD)
	addFlagContinueDirectly(applyCMD)
}
