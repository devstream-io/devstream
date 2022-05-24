package plugin

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/cmd/devstream/list"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates a new plugin.
// 1. Render template files
// 2. Persist all files
// 3. Print help information
func Create() error {
	name := viper.GetString("name")
	if name == "" {
		return fmt.Errorf("the name must be not \"\", you can specify it by --name flag")
	}
	log.Debugf("Got the name: %s.", name)

	if pluginExists(name) {
		return fmt.Errorf("Plugin name: %s is already exists", name)
	}
	log.Debugf("Got the name: %s.", name)

	p := NewPlugin(name)

	// 1. Render template files
	files, err := p.RenderTplFiles()
	if err != nil {
		log.Debugf("Failed to render template files: %s.", err)
		return err
	}
	log.Info("Render template files finished.")

	// 2. Persist all files
	if err := p.PersistFiles(files); err != nil {
		log.Debugf("Failed to persist files: %s.", err)
	}
	log.Info("Persist all files finished.")

	// 3. Print help information
	p.PrintHelpInfo()

	return nil
}

// Check whether the new name exists.
func pluginExists(name string) bool {
	pluginMp := list.PluginNamesMap()

	if _, ok := (pluginMp)[name]; ok {
		return true
	}

	return false
}
