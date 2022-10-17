package pluginengine

import (
	"fmt"
	"plugin"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

func loadPlugin(pluginDir string, tool *configmanager.Tool) (DevStreamPlugin, error) {
	mod := fmt.Sprintf("%s/%s", pluginDir, tool.GetPluginFileName())
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
