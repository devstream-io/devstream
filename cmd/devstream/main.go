package main

import (
	"github.com/spf13/cobra"
)

var (
	configFile string

	rootCMD = &cobra.Command{
		Use:   "dtm",
		Short: "DevStream is an open-source DevOps tool manager",
		Long:  `DevStream is an open-source DevOps tool manager`,
	}
)

// Execute runs the rootCMD's Execute func.
func Execute() error {
	return rootCMD.Execute()
}

func init() {
	rootCMD.PersistentFlags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")

	rootCMD.AddCommand(versionCMD)
	rootCMD.AddCommand(installCMD)
}

func main() {
	rootCMD.Execute()
}
