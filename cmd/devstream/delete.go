package main

import (
	"log"
	"os"

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
	log.Println("Delete started.")

	if err := pluginengine.Delete(configFile, continueDirectly); err != nil {
		log.Printf("Delete error: %s.", err)
		os.Exit(1)
	}

	log.Println("Delete finished.")
}
