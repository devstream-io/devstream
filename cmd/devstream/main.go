package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	configFile string

	rootCMD = &cobra.Command{
		Use:   "dtm",
		Short: "DevStream is an open-source DevOps toolchain manager",
		Long:  `DevStream is an open-source DevOps toolchain manager`,
	}
)

func init() {
	rootCMD.PersistentFlags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")

	rootCMD.AddCommand(versionCMD)
	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(installCMD)
}

func main() {
	err := rootCMD.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
