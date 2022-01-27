package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/pluginengine"
)

var (
	configFile       string
	pluginDir        string
	continueDirectly bool
	isDebug          bool

	rootCMD = &cobra.Command{
		Use:   "dtm",
		Short: `DevStream is an open-source DevOps toolchain manager`,
		Long: `DevStream is an open-source DevOps toolchain manager

######                 #####                                    
#     # ###### #    # #     # ##### #####  ######   ##   #    # 
#     # #      #    # #         #   #    # #       #  #  ##  ## 
#     # #####  #    #  #####    #   #    # #####  #    # # ## # 
#     # #      #    #       #   #   #####  #      ###### #    # 
#     # #       #  #  #     #   #   #   #  #      #    # #    # 
######  ######   ##    #####    #   #    # ###### #    # #    # 
`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCMD.PersistentFlags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
	rootCMD.PersistentFlags().StringVarP(&pluginDir, "plugin-dir", "p", pluginengine.DefaultPluginDir, "plugins directory")
	rootCMD.PersistentFlags().BoolVarP(&continueDirectly, "yes", "y", false, "apply/delete directly without confirmation")
	rootCMD.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "debug level log")

	rootCMD.AddCommand(versionCMD)
	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(applyCMD)
	rootCMD.AddCommand(deleteCMD)
	rootCMD.AddCommand(verifyCMD)
}

func initConfig() {
	viper.AutomaticEnv()
	if err := viper.BindEnv("github_token"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("dockerhub_username"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("dockerhub_token"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("trello_api_key"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindEnv("trello_token"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlags(rootCMD.Flags()); err != nil {
		log.Fatal(err)
	}
}

func initLogConf() {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		log.Infof("init... %s", logrus.GetLevel().String())
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		log.Infof("init--... %s", logrus.GetLevel().String())
	}
}

func main() {
	err := rootCMD.Execute()
	if err != nil {
		if strings.Contains(err.Error(), "unknown command \"install\"") {
			log.Fatalf("Did you mean \"dtm apply\" instead?")
		} else {
			log.Fatal(err.Error())
		}
	}
}
