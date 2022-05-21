package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var configFile string 


var applyCMD = &cobra.Command{
	Use:   "apply",
	Short: "Create or update DevOps tools according to DevStream configuration file",
	Long: `Create or update DevOps tools according to DevStream configuration file. 
DevStream will generate and execute a new plan based on the config file and the state file by default.`,
	Run: applyCMDFunc,
}

var applyConfigCMD = &cobra.Command{
	Use:   "config",
	Short: "Show configuration information",
	Long: `Show config is used for showing plugins' template configuration information.
Examples:
  dtm show config --plugin=A-PLUGIN-NAME`,
	Run: applyConfigCMDFunc,
}



func applyCMDFunc(cmd *cobra.Command, args []string) {
	log.Info("Apply started.")
	if err := pluginengine.Apply(configFile, continueDirectly); err != nil {
		log.Errorf("Apply failed => %s.", err)
		os.Exit(1)
	}
	log.Success("Apply finished.")
}


func applyconfigCMDfunc(cmd *cobra.Command, args []string){
	
//	log.Debug("Show configuration information.")
	if err := config.Apply(); err != nil {
		log.Fatal(err)
	}
}


func init() {
	
	applyCMD.AddCommand(applyCMD)
	applyCMD.PersistentFlags().StringVarP(&configFile, "config-file", "f", "config.yaml", "config file")
	viper.BindPFlag("confgFile", applyCmd.PersistentFlags().Lookup("configFile"))

	
}

func initconfig(){
	viper.AutomaticEnv()
	if err := viper.BindPFlags( applyCMD.Flags()); err != nil {
		log.Fatal(err)
	}
}

