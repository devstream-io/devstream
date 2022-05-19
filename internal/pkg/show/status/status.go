package status

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Show(configFile string) error {
	plugin := viper.GetString("plugin")
	name := viper.GetString("name")
	if name == "" {
		log.Warnf("Empty instance name. Maybe you forgot to add --name=PLUGIN_INSTANCE_NAME. The default value \"default\" will be used.")
		name = "default"
	}

	// if --plugin="" and --name="", we set the allFlag to true, it means all plugins' status need to be printed
	var allFlag = false
	if plugin == "" {
		allFlag = true
	}

	cfg, err := configloader.LoadConf(configFile)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	smgr, err := statemanager.NewManager(*cfg.State)
	if err != nil {
		log.Debugf("Failed to get State Manager: %s.", err)
		return err
	}

	if allFlag {
		return showAll(smgr)
	}
	return showOne(smgr, name, plugin)
}

// show all plugins' status
func showAll(smgr statemanager.Manager) error {
	fmt.Println()
	stateList := smgr.GetStatesMap().ToList()

	if len(stateList) == 0 {
		log.Info("No resources found.")
		return nil
	}

	var retErrs = make([]string, 0)
	for i, state := range stateList {
		fmt.Printf("================= %d/%d =================\n\n", i+1, len(stateList))
		if err := showOne(smgr, state.InstanceID, state.Name); err != nil {
			log.Errorf("Failed to show the status with <%s.%s>, error: %s.", state.InstanceID, state.Name, err)
			retErrs = append(retErrs, err.Error())
			// the "continue" here is used to tell you we don't need to return when ONE plugin show failed
			continue
		}
	}

	if len(retErrs) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(retErrs, "; "))
}

// show one plugin status
func showOne(smgr statemanager.Manager, name, plugin string) error {
	// get state from statemanager
	state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndPluginKind(name, plugin))
	if state == nil {
		return fmt.Errorf("state with (name: %s, plugin: %s) not found", name, plugin)
	}

	// get state from read
	tool := &configloader.Tool{
		InstanceID: name,
		Name:       plugin,
		DependsOn:  state.DependsOn,
		Options:    state.Options,
	}
	stateFromRead, err := pluginengine.Read(tool)
	if err != nil {
		log.Debugf("Failed to get the resource state with %s.%s. Error: %s.", name, plugin, err)
		return err
	}

	// assemble the status
	var status = &Status{}
	if reflect.DeepEqual(state.Resource, stateFromRead) {
		status.InlineStatus = state.Resource
		// set-to-nil has no effect, but make the logic more readable. Same as below.
		status.State = nil
		status.Resource = nil
	} else {
		status.InlineStatus = nil
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
