package pluginengine

import (
	"github.com/google/go-cmp/cmp"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func generateCreateAction(tool *configloader.Tool) *Change {
	return &Change{
		Tool:       tool.DeepCopy(),
		ActionName: statemanager.ActionCreate,
	}
}

func generateUpdateAction(tool *configloader.Tool) *Change {
	return &Change{
		Tool:       tool.DeepCopy(),
		ActionName: statemanager.ActionUpdate,
	}
}

func generateDeleteAction(tool *configloader.Tool) *Change {
	return &Change{
		Tool:       tool.DeepCopy(),
		ActionName: statemanager.ActionDelete,
	}
}

func drifted(r map[string]interface{}, s *statemanager.State) bool {
	return !cmp.Equal(r, s.Resource)
}

// changesForApply is to filter all the Tools in cfg that need some actions
func changesForApply(smgr statemanager.Manager, statesMap statemanager.StatesMap, cfg *configloader.Config) ([]*Change, error) {
	changes := make([]*Change, 0)
	for _, tool := range cfg.Tools {
		state := smgr.GetState(getStateKeyFromTool(&tool))

		// config not found in state
		if state == nil {
			// no need to Read resource before Create
			changes = append(changes, generateCreateAction(&tool))
			log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionCreate)
			continue
		}

		// config found in state / state isn't empty; trying to read resource status
		resource, err := Read(&tool)
		if err != nil {
			return changes, err
		}

		if resource == nil {
			// state has it but resource doesn't exist; need to create
			changes = append(changes, generateCreateAction(&tool))
			log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionCreate)
		} else if drifted(resource, state) {
			// resource drifted from state, need to update
			changes = append(changes, generateCreateAction(&tool))
			log.Infof("Change added: %s -> %s", tool.Name, statemanager.ActionUpdate)
		} else {
			// resource is the same as the state, do nothing
			log.Debugf("Tool %s state and resource are the same, not drifted, do nothing.", tool.Name)
		}

		statesMap.Delete(getStateKeyFromTool(&tool))
	}

	return changes, nil
}

// changesForDelete is to create a plan that deletes all the Tools in cfg
func changesForDelete(smgr statemanager.Manager, statesMap statemanager.StatesMap, cfg *configloader.Config) []*Change {
	changes := make([]*Change, 0)

	// reverse loop, a hack to solve dependency issues when deleting
	for i := len(cfg.Tools) - 1; i >= 0; i-- {
		tool := cfg.Tools[i]
		state := smgr.GetState(getStateKeyFromTool(&tool))
		if state == nil {
			continue
		}
		action := statemanager.ActionDelete
		changes = append(changes, &Change{
			Tool:       tool.DeepCopy(),
			ActionName: action,
		})
		log.Infof("Change added: %s -> %s", tool.Name, action)
		statesMap.Delete(tool.Name)
	}

	return changes
}
