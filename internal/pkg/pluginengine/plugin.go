package pluginengine

import (
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

// DevStreamPlugin is a struct, on which Create/Read/Update/Delete interfaces are defined.
type DevStreamPlugin interface {
	// Create, Read, and Update return two results, the first being the "state"
	Create(configmanager.RawOptions) (statemanager.ResourceStatus, error)
	Read(configmanager.RawOptions) (statemanager.ResourceStatus, error)
	Update(configmanager.RawOptions) (statemanager.ResourceStatus, error)
	// Delete returns (true, nil) if there is no error; otherwise it returns (false, error)
	Delete(configmanager.RawOptions) (bool, error)
}

// Create loads the plugin and calls the Create method of that plugin.
func Create(tool *configmanager.Tool) (statemanager.ResourceStatus, error) {
	pluginDir := viper.GetString("plugin-dir")
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	renderInstanceIDtoOptions(tool)
	return p.Create(tool.Options)
}

// Update loads the plugin and calls the Update method of that plugin.
func Update(tool *configmanager.Tool) (statemanager.ResourceStatus, error) {
	pluginDir := viper.GetString("plugin-dir")
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	renderInstanceIDtoOptions(tool)
	return p.Update(tool.Options)
}

func Read(tool *configmanager.Tool) (statemanager.ResourceStatus, error) {
	pluginDir := viper.GetString("plugin-dir")
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return nil, err
	}
	renderInstanceIDtoOptions(tool)
	return p.Read(tool.Options)
}

// Delete loads the plugin and calls the Delete method of that plugin.
func Delete(tool *configmanager.Tool) (bool, error) {
	pluginDir := viper.GetString("plugin-dir")
	p, err := loadPlugin(pluginDir, tool)
	if err != nil {
		return false, err
	}
	renderInstanceIDtoOptions(tool)
	return p.Delete(tool.Options)
}

func renderInstanceIDtoOptions(tool *configmanager.Tool) {
	tool.Options["instanceID"] = interface{}(tool.InstanceID)
}
