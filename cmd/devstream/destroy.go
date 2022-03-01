package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/pluginengine"
	"github.com/merico-dev/stream/pkg/util/log"
)

var destroyCMD = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy DevOps tools deployment according to DevStream configuration file & state file",
	Long:  `Destroy DevOps tools deployment according to DevStream configuration file & state file.`,
	Run:   destroyCMDFunc,
}

func destroyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Destroy started.")
	if err := pluginengine.Destroy(continueDirectly); err != nil {
		log.Errorf("Destroy failed => %s.", err)
		os.Exit(1)
	}
	log.Success("Destroy finished.")
}
