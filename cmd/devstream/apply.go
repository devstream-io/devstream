package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/pluginengine"
)

var applyCMD = &cobra.Command{
	Use:   "apply",
	Short: "Create or update DevOps tools according to DevStream configuration file",
	Long: `Create or update DevOps tools according to DevStream configuration file. 
DevStream will generate and execute a new plan based on the config file and the state file by default.`,
	Run: applyCMDFunc,
}

func applyCMDFunc(cmd *cobra.Command, args []string) {
	log.Println("Apply started.")

	err := pluginengine.Apply(configFile)
	if err != nil {
		log.Printf("Apply error: %s. Maybe you forgot to run \"dtm init\" first?", err)
		os.Exit(1)
	}

	log.Println("Apply finished.")
}
