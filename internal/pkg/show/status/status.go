package status

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Show() error {
	plugin := viper.GetString("plugin")
	name := viper.GetString("name")

	if plugin == "" && name != "" {
		return fmt.Errorf("empty plugin name. Maybe you forgot to add --plugin=PLUGIN_NAME")
	}
	if name == "" && plugin != "" {
		return fmt.Errorf("empty instance name. Maybe you forgot to add --name=PLUGIN_INSTANCE_NAME")
	}

	var allFlag = false
	if name == "" && plugin == "" {
		allFlag = true
	}

	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Debugf("Failed to get State Manager: %s.", err)
		return err
	}

	if allFlag {
		return showAll(smgr)
	}
	return showOne(smgr, name, plugin)
}

func showAll(smgr statemanager.Manager) error {
	return nil
}

func showOne(smgr statemanager.Manager, name, plugin string) error {
	// get state from statemanager
	state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndPluginKind(name, plugin))
	if state == nil {
		return fmt.Errorf("state with (name: %s, plugin: %s) not found", name, plugin)
	}

	// get state from read
	tool := &configloader.Tool{
		Name:      name,
		Plugin:    plugin,
		DependsOn: state.DependsOn,
		Options:   state.Options,
	}
	stateFromRead, err := pluginengine.Read(tool)
	if err != nil {
		log.Debugf("Failed to get the resource state with %s.%s. Error: %s.", name, plugin, err)
		return err
	}

	// assemble the status
	var status = &Status{}
	if reflect.DeepEqual(state.Resource, stateFromRead) {
		status.InnerStatus = state.Resource
		// set-to-nil has no effect, but make the logic more readable. Same as below.
		status.State = nil
		status.Resource = nil
	} else {
		status.InnerStatus = nil
		status.State = state.Resource
		status.Resource = stateFromRead
	}

	// get the output
	output, err := NewOutput(name, plugin, state.Options, status)
	if err != nil {
		log.Debugf("Failed to get the output: %s.", err)
	}

	// print
	return output.Print()
}
