package pluginengine

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

func generateCreateAction(tool *configloader.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionCreate, description)
}

func generateUpdateAction(tool *configloader.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionUpdate, description)
}

func generateDeleteAction(tool *configloader.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionDelete, description)
}

func generateDeleteActionFromState(state statemanager.State) *Change {
	return &Change{
		Tool: &configloader.Tool{
			Name:    state.Name,
			Plugin:  state.Plugin,
			Options: state.Options,
		},
		ActionName: statemanager.ActionDelete,
	}
}

func generateAction(tool *configloader.Tool, action statemanager.ComponentAction, description string) *Change {
	return &Change{
		Tool:        tool.DeepCopy(),
		ActionName:  action,
		Description: description,
	}
}

func drifted(a, b map[string]interface{}) bool {
	// nil vs empty map
	if cmp.Equal(a, b, cmpopts.EquateEmpty()) {
		return false
	}

	log.Debug(cmp.Diff(a, b))
	return !cmp.Equal(a, b)
}

// changesForApply generates "changes" according to:
// - config
// - state
// - resource status (by calling the Read() interface of the plugin)
func changesForApply(smgr statemanager.Manager, cfg *configloader.Config) ([]*Change, error) {
	changes := make([]*Change, 0)

	// 1, create a temporary state map used to store unprocessed tools.
	tmpStatesMap := smgr.GetStatesMap().DeepCopy()

	// 2, for each tool in the config, generate changes.
	for _, tool := range cfg.Tools {
		state := smgr.GetState(getStateKeyFromTool(&tool))

		if state == nil {
			// tool not in the state, create, no need to Read resource before Create
			description := fmt.Sprintf("Tool < %s > found in config but doesn't exist in the state, will be created.", tool.Name)
			changes = append(changes, generateCreateAction(&tool, description))
		} else {
			// tool found in the state
			if drifted(tool.Options, state.Options) {
				// tool's config differs from State's, Update
				description := fmt.Sprintf("Tool < %s > config drifted from the state, will be updated.", tool.Name)
				changes = append(changes, generateUpdateAction(&tool, description))
			} else {
				// tool's config is the same as State's

				// read resource status
				resource, err := Read(&tool)
				if err != nil {
					return changes, err
				}

				if resource == nil {
					// tool exists in the state, but resource doesn't exist, Create
					description := fmt.Sprintf("Tool < %s > state found but it seems the tool isn't created, will be created.", tool.Name)
					changes = append(changes, generateCreateAction(&tool, description))
				} else if drifted(resource, state.Resource) {
					// resource drifted from state, Update
					description := fmt.Sprintf("Tool < %s > drifted from the state, will be updated.", tool.Name)
					changes = append(changes, generateUpdateAction(&tool, description))
				} else {
					// resource is the same as the state, do nothing
				}
			}
		}

		// delete the tool from the temporary state map since it's already been processed above
		tmpStatesMap.Delete(getStateKeyFromTool(&tool))
	}

	// what's left in the temporary state map "tmpStatesMap" contains tools that:
	// - have a state (probably created previously)
	// - don't have a definition in the config (probably deleted by the user)
	// thus, we need to generate a "delete" change for it.
	tmpStatesMap.Range(func(key, value interface{}) bool {
		changes = append(changes, generateDeleteActionFromState(value.(statemanager.State)))
		log.Infof("Change added: %s -> %s", key.(string), statemanager.ActionDelete)
		return true
	})

	return changes, nil
}

// changesForDelete is to create a plan that deletes all the Tools in cfg
func changesForDelete(smgr statemanager.Manager, cfg *configloader.Config) []*Change {
	changes := make([]*Change, 0)
	tmpStates := smgr.GetStatesMap().DeepCopy()

	// reverse loop, a hack to solve dependency issues when deleting
	for i := len(cfg.Tools) - 1; i >= 0; i-- {
		tool := cfg.Tools[i]
		state := smgr.GetState(getStateKeyFromTool(&tool))
		if state == nil {
			continue
		}
		description := fmt.Sprintf("Tool < %s > will be deleted.", tool.Name)
		changes = append(changes, generateDeleteAction(&tool, description))
		tmpStates.Delete(tool.Name)
	}

	return changes
}
