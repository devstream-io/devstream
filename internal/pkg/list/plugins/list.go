package plugins

import (
	"fmt"
	"strings"
	// "github.com/spf13/viper"
	// "github.com/devstream-io/devstream/pkg/util/log"
)

// List all of plugins name
func List(PluginsName string) error {
	listPluginsName := strings.Fields(PluginsName)
	for _, pluginName := range listPluginsName {
		fmt.Println(pluginName)
	}
	return nil
}
