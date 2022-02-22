package pluginengine

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// Change is a wrapper with a single Tool and its Action should be execute.
type Change struct {
	Tool       *configloader.Tool
	ActionName statemanager.ComponentAction
	Result     *ChangeResult
}

// ChangeResult holds the result with a change action.
type ChangeResult struct {
	Succeeded   bool
	Error       error
	Time        string
	ReturnValue map[string]interface{}
}

func (c *Change) String() string {
	return fmt.Sprintf("\n{\n  ActionName: %s,\n  Tool: {Name: %s, Plugin: {Kind: %s, Version: %s}}\n}",
		c.ActionName, c.Tool.Name, c.Tool.Plugin.Kind, c.Tool.Plugin.Version)
}

// GetChangesForApply takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be execute is included in this Plan.changes.
func GetChangesForApply(smgr statemanager.Manager, cfg *configloader.Config) ([]*Change, error) {
	if cfg == nil {
		return make([]*Change, 0), nil
	}

	data, err := smgr.Read()
	if err == nil {
		statesMap := statemanager.NewStatesMap()
		tmpMap := make(map[string]statemanager.State)
		if err := yaml.Unmarshal(data, tmpMap); err != nil {
			log.Errorf("Failed to unmarshal the state file < %s >. error: %s", local.DefaultStateFile, err)
			log.Errorf("Reading the state file failed, it might have been compromised/modified by someone other than DTM.")
			log.Errorf("The state file is managed by DTM automatically. Please do not modify it yourself.")
			return make([]*Change, 0), err
		}
		for k, v := range tmpMap {
			statesMap.Store(k, v)
		}
		smgr.SetStatesMap(statesMap)

		log.Success("Succeeded initializing StatesMap.")
	} else {
		log.Errorf("Failed to initialize StatesMap. Error: (%s). Try to initialize the StatesMap.", err)
	}

	// calculate changes from config and state
	tmpStates := smgr.GetStatesMap().DeepCopy()
	changes, err := changesForApply(smgr, tmpStates, cfg)
	if err != nil {
		return nil, err
	}

	log.Debugf("Changes for the plan:")
	for _, c := range changes {
		log.Debugf(c.String())
	}

	return changes, nil
}

// GetChangesForDelete takes "State Manager" & "Config" then do some calculation and return a Plan to delete all plugins in the Config.
// All actions should be execute is included in this Plan.changes.
func GetChangesForDelete(smgr statemanager.Manager, cfg *configloader.Config) ([]*Change, error) {
	if cfg == nil {
		return make([]*Change, 0), nil
	}

	data, err := smgr.Read()
	// TODO(ironcore864): duplicated code; needs to be refactored.
	if err == nil {
		statesMap := statemanager.NewStatesMap()
		tmpMap := make(map[string]statemanager.State)
		if err := yaml.Unmarshal(data, tmpMap); err != nil {
			log.Fatalf("Devstream.statesMap format error.")
			return make([]*Change, 0), err
		}
		for k, v := range tmpMap {
			statesMap.Store(k, v)
		}
		smgr.SetStatesMap(statesMap)
		log.Success("Succeeded initializing StatesMap.")
	} else {
		log.Errorf("Failed to initialize StatesMap. Error: (%s). Try to initialize the StatesMap.", err)
	}

	tmpStates := smgr.GetStatesMap().DeepCopy()
	changes := changesForDelete(smgr, tmpStates, cfg)

	log.Debugf("Changes for the plan:")
	for _, c := range changes {
		log.Debugf(c.String())
	}

	return changes, nil
}

func execute(smgr statemanager.Manager, changes []*Change) map[string]error {
	errorsMap := make(map[string]error)

	log.Info("Start executing the plan.")
	numOfChanges := len(changes)
	log.Infof("Changes count: %d.", numOfChanges)

	for i, c := range changes {
		log.Separatorf("Processing progress: %d/%d.", i+1, numOfChanges)
		log.Infof("Processing: %s -> %s ...", c.Tool.Name, c.ActionName)

		var succeeded bool
		var err error
		var returnValue map[string]interface{}

		switch c.ActionName {
		case statemanager.ActionCreate:
			if returnValue, err = Create(c.Tool); err == nil {
				succeeded = true
			}
		case statemanager.ActionUpdate:
			if returnValue, err = Update(c.Tool); err == nil {
				succeeded = true
			}
		case statemanager.ActionDelete:
			succeeded, err = Delete(c.Tool)
		}

		if err != nil {
			key := fmt.Sprintf("%s-%s", c.Tool.Name, c.ActionName)
			errorsMap[key] = err
		}

		c.Result = &ChangeResult{
			Succeeded:   succeeded,
			Error:       err,
			Time:        time.Now().Format(time.RFC3339),
			ReturnValue: returnValue,
		}

		err = handleResult(smgr, c)
		if err != nil {
			errorsMap["handle-result"] = err
		}
	}

	return errorsMap
}

// handleResult is used to Write the latest StatesMap to the Backend.
func handleResult(smgr statemanager.Manager, change *Change) error {
	log.Debugf("Start: \n%s", string(smgr.GetStatesMap().Format()))
	defer func() {
		log.Debugf("End:\n%s", string(smgr.GetStatesMap().Format()))
	}()

	if !change.Result.Succeeded {
		log.Errorf("Plugin %s %s failed.", change.Tool.Name, change.ActionName)
		return fmt.Errorf("plugin %s %s failed", change.Tool.Name, change.ActionName)
	}

	if change.ActionName == statemanager.ActionDelete {
		key := getStateKeyFromTool(change.Tool)
		log.Infof("Prepare to delete '%s' from States", key)
		smgr.DeleteState(key)
		log.Successf("Plugin %s delete done.", change.Tool.Name)
		return smgr.Write(smgr.GetStatesMap().Format())
	}

	key := getStateKeyFromTool(change.Tool)
	state := statemanager.State{
		Name:     change.Tool.Name,
		Plugin:   change.Tool.Plugin,
		Resource: change.Result.ReturnValue,
	}
	smgr.AddState(key, state)

	log.Successf("Plugin %s %s done.", change.Tool.Name, change.ActionName)
	return smgr.Write(smgr.GetStatesMap().Format())
}
