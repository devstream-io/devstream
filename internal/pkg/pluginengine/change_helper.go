package pluginengine

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
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
			InstanceID: state.InstanceID,
			Name:       state.Name,
			Options:    state.Options,
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

	// 1. create a temporary state map used to store unprocessed tools.
	tmpStatesMap := smgr.GetStatesMap().DeepCopy()

	// 2. handle dependency and sort the tools in the config into "batches" of tools
	var batchesOfTools [][]configloader.Tool
	// the elements in batchesOfTools are sorted "batches"
	// and each element/batch is a list of tools that, in theory, can run in parallel
	// that is to say, the tools in the same batch won't depend on each other
	batchesOfTools, err := topologicalSort(cfg.Tools)
	if err != nil {
		return changes, err
	}

	// 3. generate changes for each tool
	for _, batch := range batchesOfTools {
		for _, tool := range batch {
			state := smgr.GetState(statemanager.StateKeyGenerateFunc(&tool))

			if state == nil {
				// tool not in the state, create, no need to Read resource before Create
				description := fmt.Sprintf("Tool (%s/%s) found in config but doesn't exist in the state, will be created.", tool.Name, tool.InstanceID)
				changes = append(changes, generateCreateAction(&tool, description))
			} else {
				// tool found in the state

				// first, handle possible "outputs" referneces in the tool's config
				// ignoring errors, since at this stage we are calculating changes, and the dependency might not have its output in the state yet
				_ = HandleOutputsReferences(smgr, tool.Options)

				if drifted(tool.Options, state.Options) {
					// tool's config differs from State's, Update
					description := fmt.Sprintf("Tool (%s/%s) config drifted from the state, will be updated.", tool.Name, tool.InstanceID)
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
						description := fmt.Sprintf("Tool (%s/%s) state found but it seems the tool isn't created, will be created.", tool.Name, tool.InstanceID)
						changes = append(changes, generateCreateAction(&tool, description))
					} else if drifted(resource, state.Resource) {
						// resource drifted from state, Update
						description := fmt.Sprintf("Tool (%s/%s) drifted from the state, will be updated.", tool.Name, tool.InstanceID)
						changes = append(changes, generateUpdateAction(&tool, description))
					} else {
						// resource is the same as the state, do nothing
						log.Debugf("Tool (%s/%s) is the same as the state, do nothing.", tool.Name, tool.InstanceID)
					}
				}
			}

			// delete the tool from the temporary state map since it's already been processed above
			tmpStatesMap.Delete(statemanager.StateKeyGenerateFunc(&tool))
		}
	}

	// what's left in the temporary state map "tmpStatesMap" contains tools that:
	// - have a state (probably created previously)
	// - don't have a definition in the config (probably deleted by the user)
	// thus, we need to generate a "delete" change for it.
	tmpStatesMap.Range(func(key, value interface{}) bool {
		changes = append(changes, generateDeleteActionFromState(value.(statemanager.State)))
		log.Infof("Change added: %s -> %s", key.(statemanager.StateKey), statemanager.ActionDelete)
		return true
	})

	return changes, nil
}

// changesForDelete is to create a plan that deletes all the Tools in cfg
func changesForDelete(smgr statemanager.Manager, cfg *configloader.Config, isForceDelete bool) ([]*Change, error) {
	changes := make([]*Change, 0)
	batchesOfTools, err := topologicalSort(cfg.Tools)
	if err != nil {
		return changes, err
	}

	for i := len(batchesOfTools) - 1; i >= 0; i-- {
		batch := batchesOfTools[i]
		for _, tool := range batch {
			if !isForceDelete {
				state := smgr.GetState(statemanager.StateKeyGenerateFunc(&tool))
				if state == nil {
					continue
				}
			}

			description := fmt.Sprintf("Tool (%s/%s) will be deleted.", tool.Name, tool.InstanceID)
			changes = append(changes, generateDeleteAction(&tool, description))
		}
	}

	return changes, nil
}

func GetChangesForDestroy(smgr statemanager.Manager) ([]*Change, error) {
	changes := make([]*Change, 0)

	// rebuilding tools from config
	// because destroy will only be used when the config file is missing
	var tools []configloader.Tool
	for _, state := range smgr.GetStatesMap().ToList() {
		tool := configloader.Tool{
			InstanceID: state.InstanceID,
			Name:       state.Name,
			DependsOn:  state.DependsOn,
			Options:    state.Options,
		}
		tools = append(tools, tool)
	}

	batchesOfTools, err := topologicalSort(tools)
	if err != nil {
		return changes, err
	}

	// reverse, for deletion
	for i := len(batchesOfTools) - 1; i >= 0; i-- {
		batch := batchesOfTools[i]
		for _, tool := range batch {
			description := fmt.Sprintf("Tool (%s/%s) will be deleted.", tool.Name, tool.InstanceID)
			changes = append(changes, generateDeleteAction(&tool, description))
		}
	}

	return changes, nil
}
