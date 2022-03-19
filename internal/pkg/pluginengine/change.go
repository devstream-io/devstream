package pluginengine

import (
	"fmt"
	"time"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Change is a wrapper with a single Tool and its Action should be execute.
type Change struct {
	Tool        *configloader.Tool
	ActionName  statemanager.ComponentAction
	Result      *ChangeResult
	Description string
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

type CommandType string

const (
	CommandApply  CommandType = "apply"
	CommandDelete CommandType = "delete"
)

// GetChangesForApply takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be execute is included in this Plan.changes.
func GetChangesForApply(smgr statemanager.Manager, cfg *configloader.Config) ([]*Change, error) {
	return getChanges(smgr, cfg, CommandApply, false)
}

// GetChangesForDelete takes "State Manager" & "Config" then do some calculation and return a Plan to delete all plugins in the Config.
// All actions should be execute is included in this Plan.changes.
func GetChangesForDelete(smgr statemanager.Manager, cfg *configloader.Config, isForceDelete bool) ([]*Change, error) {
	return getChanges(smgr, cfg, CommandDelete, isForceDelete)
}

func getChanges(smgr statemanager.Manager, cfg *configloader.Config, commandType CommandType, isForceDelete bool) ([]*Change, error) {
	if cfg == nil {
		return make([]*Change, 0), nil
	}
	log.Debug("isForce:", isForceDelete)
	// calculate changes from config and state
	var changes []*Change
	var err error
	if commandType == CommandApply {
		changes, err = changesForApply(smgr, cfg)
	} else if commandType == CommandDelete {
		if isForceDelete {
			changes = changesForForceDelete(smgr, cfg)
		} else {
			changes = changesForDelete(smgr, cfg)
		}
	} else {
		log.Fatalf("That's impossible!")
	}

	if err != nil {
		return nil, err
	}

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
		log.Infof("Processing: %s(kind: %s) -> %s ...", c.Tool.Name, c.Tool.Plugin.Kind, c.ActionName)

		var succeeded bool
		var err error
		var returnValue map[string]interface{}

		log.Debugf("Tool's raw changes are: %s.", c.Tool.Options)

		errs := HandleOutputsReferences(smgr, c.Tool.Options)
		if len(errs) != 0 {
			succeeded = false

			for _, e := range errs {
				log.Errorf("Error: %s", e)
			}
			log.Errorf("The outputs reference in tool %s (%s) can't be resolved. Please double check your config.", c.Tool.Name, c.Tool.Plugin.Kind)

			// not executing this change since its input isn't valid
			continue
		}

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
			key := fmt.Sprintf("%s/%s-%s", c.Tool.Plugin.Kind, c.Tool.Name, c.ActionName)
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
	log.Separatorf("Processing done.")

	return errorsMap
}

// handleResult is used to Write the latest StatesMap to the Backend.
func handleResult(smgr statemanager.Manager, change *Change) error {
	log.Debugf("Start: \n%s", string(smgr.GetStatesMap().Format()))
	defer func() {
		log.Debugf("End:\n%s", string(smgr.GetStatesMap().Format()))
	}()

	if !change.Result.Succeeded {
		// do nothing when the change failed
		return nil
	}

	if change.ActionName == statemanager.ActionDelete {
		key := statemanager.StateKeyGenerateFunc(change.Tool)
		log.Infof("Prepare to delete '%s' from States.", key)
		err := smgr.DeleteState(key)
		if err != nil {
			log.Debugf("Failed to delete state %s: %s.", key, err)
			return err
		}
		log.Successf("Plugin %s/%s delete done.", change.Tool.Name, change.Tool.Plugin.Kind)
		return nil
	}

	key := statemanager.StateKeyGenerateFunc(change.Tool)
	state := statemanager.State{
		Name:     change.Tool.Name,
		Plugin:   change.Tool.Plugin,
		Options:  change.Tool.Options,
		Resource: change.Result.ReturnValue,
	}
	err := smgr.AddState(key, state)
	if err != nil {
		log.Debugf("Failed to add state %s: %s.", key, err)
		return err
	}
	log.Successf("Plugin %s(%s) %s done.", change.Tool.Name, change.Tool.Plugin.Kind, change.ActionName)
	return nil
}
