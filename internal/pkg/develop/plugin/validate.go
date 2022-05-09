package plugin

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Validate() error {
	validateAll := viper.GetBool("all")

	if validateAll {
		log.Info("Start validating all plugins.\n")
		return ValidatePlugins()
	}

	name := viper.GetString("name")
	if name == "" {
		return fmt.Errorf("the name must be not \"\", you can specify it by --name flag")
	}
	log.Debugf("Got the name: %s.", name)

	return ValidatePlugin(name)
}

// Validate all plugins
// calling ValidatePlugin() via all plugins name
func ValidatePlugins() error {
	listPluginsName := list.PluginsNameSlice()

	for _, pluginName := range listPluginsName {
		log.Infof("===== start validating <%s> =====", pluginName)
		if err := ValidatePlugin(pluginName); err != nil {
			return err
		}
		log.Infof("===== Finished validate <%s> =====\n", pluginName)
	}

	return nil
}

// Validate a plugin.
// 1. Render template files
// 2. Validate need validate files
func ValidatePlugin(pluginName string) error {
	p := NewPlugin(pluginName)

	// 1. Render template files
	files, err := p.RenderTplFiles()
	if err != nil {
		log.Debugf("Failed to render template files: %s.", err)
		return err
	}
	log.Debug("Render template files finished.")

	// 2. Validate need validate files
	if err := p.ValidateFiles(files); err != nil {
		log.Debugf("Failed to validate files: %s.", err)
		return err
	}

	return nil
}
