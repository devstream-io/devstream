package main

import (
	"os"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/pluginengine"
)

var deleteCMD = &cobra.Command{
	Use:   "delete",
	Short: "Delete DevOps tools according to DevStream configuration file",
	Long: `Delete DevOps tools according to DevStream configuration file. 
DevStream will delete everything defined in the config file, regardless of the state.`,
	Run: deleteCMDFunc,
}

func deleteCMDFunc(cmd *cobra.Command, args []string) {
	initLogConf()
	log.Info("Delete started.")

	if err := pluginengine.Delete(configFile, continueDirectly); err != nil {
		log.Errorf("Delete error: %s.", err)
		os.Exit(1)
	}

	log.Success("Delete finished.")
}
