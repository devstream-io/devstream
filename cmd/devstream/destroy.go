package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var destroyCMD = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy DevOps tools deployment according to DevStream configuration file & state file",
	Long:  `Destroy DevOps tools deployment according to DevStream configuration file & state file.`,
	Run:   destroyCMDFunc,
}

func destroyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Destroy started.")
	if err := pluginengine.Destroy(configFile, continueDirectly); err != nil {
		log.Errorf("Destroy failed => %s.", err)
		os.Exit(1)
	}
	log.Success("Destroy finished.")
}

func init() {
	destroyCMD.Flags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
	destroyCMD.Flags().BoolVarP(&continueDirectly, "yes", "y", false, "destroy directly without confirmation")
}
