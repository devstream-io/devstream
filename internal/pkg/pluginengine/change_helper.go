package pluginengine

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

func generateCreateAction(tool *configloader.Tool) *Change {
	return generateAction(tool, statemanager.ActionCreate)
}

func generateUpdateAction(tool *configloader.Tool) *Change {
	return generateAction(tool, statemanager.ActionUpdate)
}

func generateDeleteAction(tool *configloader.Tool) *Change {
	return generateAction(tool, statemanager.ActionDelete)
}

func generateAction(tool *configloader.Tool, action statemanager.ComponentAction) *Change {
	return &Change{
		Tool:       tool.DeepCopy(),
		ActionName: action,
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

// changesForApply is to filter all the Tools in cfg that need some actions
func changesForApply(smgr statemanager.Manager, cfg *configloader.Config) ([]*Change, error) {
	changes := make([]*Change, 0)
	tmpStates := smgr.GetStatesMap().DeepCopy()

	for _, tool := range cfg.Tools {
		state := smgr.GetState(getStateKeyFromTool(&tool))

		if state == nil {
			// tool not in state, create, no need to Read resource before Create
			changes = append(changes, generateCreateAction(&tool))
			log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionCreate)
		} else {
			// tool found in state
			if drifted(tool.Options, state.Options) {
				log.Debugf("Tool %s %s config options drifted from state.", tool.Name, tool.Plugin.Kind)
				// tool's config differs from State's, Update
				changes = append(changes, generateUpdateAction(&tool))
				log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionUpdate)
			} else {
				// tool's config is the same as State's

				// read resource status
				resource, err := Read(&tool)
				if err != nil {
					return changes, err
				}

				if resource == nil {
					// tool exists in state, but resource doesn't exist, Create
					changes = append(changes, generateCreateAction(&tool))
					log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionCreate)
				} else if drifted(resource, state.Resource) {
					log.Debugf("Tool %s %s resource drifted from state.", tool.Name, tool.Plugin.Kind)
					// resource drifted from state, Update
					changes = append(changes, generateUpdateAction(&tool))
					log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionUpdate)
				} else {
					// resource is the same as the state, do nothing
					log.Debugf("Tool %s state and resource are the same, not drifted, do nothing.", tool.Name)
				}
			}
		}

		tmpStates.Delete(getStateKeyFromTool(&tool))
	}

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

		changes = append(changes, generateDeleteAction(&tool))
		log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionDelete)
		tmpStates.Delete(tool.Name)
	}

	return changes
}
