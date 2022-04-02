package pluginengine

import (
	"fmt"
	"plugin"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

func getPluginDir() string {
	var pluginDir string
	if pluginDir = viper.GetString("plugin-dir"); pluginDir == "" {
		pluginDir = DefaultPluginDir
	}
	return pluginDir
}

func loadPlugin(pluginDir string, tool *configloader.Tool) (DevStreamPlugin, error) {
	mod := fmt.Sprintf("%s/%s", pluginDir, configloader.GetPluginFileName(tool))
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}

	var devStreamPlugin DevStreamPlugin
	symDevStreamPlugin, err := plug.Lookup("DevStreamPlugin")
	if err != nil {
		return nil, err
	}

	devStreamPlugin, ok := symDevStreamPlugin.(DevStreamPlugin)
	if !ok {
		return nil, fmt.Errorf("DevStreamPlugin type error")
	}

	return devStreamPlugin, nil
}
