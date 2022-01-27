package main

import (
	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/merico-dev/stream/internal/pkg/pluginengine"
)

var verifyCMD = &cobra.Command{
	Use:   "verify",
	Short: "Verify DevOps tools according to DevStream configuration file",
	Long:  `Verify DevOps tools according to DevStream configuration file.`,
	Run:   verifyCMDFunc,
}

func verifyCMDFunc(cmd *cobra.Command, args []string) {
	initLogConf()
	log.Info("Verify started.")
	healthy, err := pluginengine.CheckHealthy(configFile)
	if err != nil {
		log.Fatalf("Verify error: %s.", err)
	}

	if healthy {
		log.Success("all tools are healthy")
	} else {
		log.Error("some tools are NOT healthy!!!")
	}
}
