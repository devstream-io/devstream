package main

import (
	"github.com/spf13/cobra"

	"github.com/devstream-io/devstream/internal/pkg/develop"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	name string
	all  bool
)

var developCMD = &cobra.Command{
	Use:   "develop",
	Short: "Develop is used for develop a new plugin",
}

var developCreatePluginCMD = &cobra.Command{
	Use:   "create-plugin",
	Short: "Create a new plugin",
	Long: `Create-plugin is used for creating a new plugin.
Exampls:
  dtm develop create-plugin --name=YOUR-PLUGIN-NAME`,
	Run: developCreateCMDFunc,
}

var developValidatePluginCMD = &cobra.Command{
	Use:   "validate-plugin",
	Short: "Validate a plugin",
	Long: `Validate-plugin is used for validating an existing plugin or all plugins.
Examples:
  dtm develop validate-plugin --name=YOUR-PLUGIN-NAME,
  dtm develop validate-plugin --all`,
	Run: developValidateCMDFunc,
}

func developCreateCMDFunc(cmd *cobra.Command, args []string) {
	if err := develop.CreatePlugin(); err != nil {
		log.Fatal(err)
	}
}

func developValidateCMDFunc(cmd *cobra.Command, args []string) {
	if err := develop.ValidatePlugin(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	developCMD.AddCommand(developCreatePluginCMD)
	developCMD.AddCommand(developValidatePluginCMD)

	developCreatePluginCMD.PersistentFlags().StringVarP(&name, "name", "n", "", "specify name with the new plugin")
	developValidatePluginCMD.PersistentFlags().StringVarP(&name, "name", "n", "", "specify name with the new plugin")
	developValidatePluginCMD.PersistentFlags().BoolVarP(&all, "all", "a", false, "validate all plugins")
}
