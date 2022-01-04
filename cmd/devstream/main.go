package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	pluginDir  string

	rootCMD = &cobra.Command{
		Use:   "dtm",
		Short: "DevStream is an open-source DevOps toolchain manager",
		Long:  `DevStream is an open-source DevOps toolchain manager`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCMD.PersistentFlags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
	rootCMD.PersistentFlags().StringVarP(&pluginDir, "plugin-dir", "p", ".devstream/plugins", "plugins directory")

	rootCMD.AddCommand(versionCMD)
	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(applyCMD)
}

func initConfig() {
	viper.AutomaticEnv()
	if err := viper.BindEnv("github_token"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlags(rootCMD.Flags()); err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := rootCMD.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
