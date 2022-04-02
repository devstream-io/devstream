package pluginengine

import (
	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

const DefaultPluginDir = ".devstream"

// DevStreamPlugin is a struct, on which Create/Read/Update/Delete interfaces are defined.
type DevStreamPlugin interface {
	// Create, Read, and Update return two results, the first being the "state"
	Create(map[string]interface{}) (map[string]interface{}, error)
	Read(map[string]interface{}) (map[string]interface{}, error)
	Update(map[string]interface{}) (map[string]interface{}, error)
	// Delete returns (true, nil) if there is no error; otherwise it returns (false, error)
	Delete(map[string]interface{}) (bool, error)
}

// Create loads the plugin and calls the Create method of that plugin.
func Create(tool *configloader.Tool) (map[string]interface{}, error) {
	pluginDir := getPluginDir()
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	return p.Create(tool.Options)
}

// Update loads the plugin and calls the Update method of that plugin.
func Update(tool *configloader.Tool) (map[string]interface{}, error) {
	pluginDir := getPluginDir()
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	return p.Update(tool.Options)
}

func Read(tool *configloader.Tool) (map[string]interface{}, error) {
	pluginDir := getPluginDir()
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	return p.Read(tool.Options)
}

// Delete loads the plugin and calls the Delete method of that plugin.
func Delete(tool *configloader.Tool) (bool, error) {
	pluginDir := getPluginDir()
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return false, err
	}
	return p.Delete(tool.Options)
}
