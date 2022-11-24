package pluginengine

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func generateCreateAction(tool *configmanager.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionCreate, description)
}

func generateUpdateAction(tool *configmanager.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionUpdate, description)
}

func generateDeleteAction(tool *configmanager.Tool, description string) *Change {
	return generateAction(tool, statemanager.ActionDelete, description)
}

func generateDeleteActionFromState(state statemanager.State) *Change {
	return &Change{
		Tool: &configmanager.Tool{
			InstanceID: state.InstanceID,
			Name:       state.Name,
			Options:    state.Options,
		},
		ActionName: statemanager.ActionDelete,
	}
}

func generateAction(tool *configmanager.Tool, action statemanager.ComponentAction, description string) *Change {
	return &Change{
		Tool:        tool.DeepCopy(),
		ActionName:  action,
		Description: description,
	}
}

func ResourceDrifted(resourceStatusFromState, resourceStatusFromRead statemanager.ResourceStatus) (bool, error) {
	// nil vs empty map
	if cmp.Equal(resourceStatusFromState, resourceStatusFromRead, cmpopts.EquateEmpty()) {
		return false, nil
	}

	// use marshal and unmarshal to remove type details
	resourceStatusFromStateBytes, err := yaml.Marshal(resourceStatusFromState)
	if err != nil {
		return true, err
	}
	resourceStatusFromStateWithoutType := map[string]interface{}{}
	err = yaml.Unmarshal(resourceStatusFromStateBytes, &resourceStatusFromStateWithoutType)
	if err != nil {
		return true, err
	}

	resourceStatusFromReadBytes, err := yaml.Marshal(resourceStatusFromRead)
	if err != nil {
		return true, err
	}
	resourceStatusFromReadWithoutType := map[string]interface{}{}
	err = yaml.Unmarshal(resourceStatusFromReadBytes, &resourceStatusFromReadWithoutType)
	if err != nil {
		return true, err
	}

	log.Debug(cmp.Diff(resourceStatusFromStateWithoutType, resourceStatusFromReadWithoutType))

	return !cmp.Equal(resourceStatusFromStateWithoutType, resourceStatusFromReadWithoutType), nil
}

func drifted(a, b map[string]interface{}) bool {
	// nil vs empty map
	if cmp.Equal(a, b, cmpopts.EquateEmpty()) {
		return false
	}

	log.Debugf("detect tool sate changed => %s", cmp.Diff(a, b))
	return !cmp.Equal(a, b)
}

// GetChangesForApply takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be executed is included in this Plan.changes.
// It generates "changes" according to:
// - config
// - state
// - resource status (by calling the Read() interface of the plugin)
func GetChangesForApply(smgr statemanager.Manager, cfg *configmanager.Config) (changes []*Change, err error) {
	// 1. create a temporary state map used to store unprocessed tools.
	tmpStatesMap := smgr.GetStatesMap().DeepCopy()

	// 2. handle dependency and sort the tools in the config into "batches" of tools
	var batchesOfTools []configmanager.Tools
	// the elements in batchesOfTools are sorted "batches"
	// and each element/batch is a list of tools that, in theory, can run in parallel
	// that is to say, the tools in the same batch won't depend on each other
	batchesOfTools, err = topologicalSort(cfg.Tools)
	if err != nil {
		return changes, err
	}

	// 3. generate changes for each tool
	for _, batch := range batchesOfTools {
		for _, tool := range batch {
			state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndInstanceID(tool.Name, tool.InstanceID))

			if state == nil {
				// tool not in the state, create, no need to Read resource before Create
				description := fmt.Sprintf("Tool (%s/%s) found in config but doesn't exist in the state, will be created.", tool.Name, tool.InstanceID)
				changes = append(changes, generateCreateAction(tool, description))
			} else {
				// tool found in the state

				// first, handle possible "outputs" references in the tool's config
				// ignoring errors, since at this stage we are calculating changes, and the dependency might not have its output in the state yet
				_ = HandleOutputsReferences(smgr, tool.Options)

				if drifted(state.Options, tool.Options) {
					// tool's config differs from State's, Update
					description := fmt.Sprintf("Tool (%s/%s) config drifted from the state, will be updated.", tool.Name, tool.InstanceID)
					changes = append(changes, generateUpdateAction(tool, description))
				} else {
					// tool's config is the same as State's

					// read resource status
					resource, err := Read(tool)
					if err != nil {
						return changes, err
					}

					if resource == nil {
						// tool exists in the state, but resource doesn't exist, Create
						description := fmt.Sprintf("Tool (%s/%s) state found but it seems the tool isn't created, will be created.", tool.Name, tool.InstanceID)
						changes = append(changes, generateCreateAction(tool, description))
					} else if drifted, err := ResourceDrifted(state.ResourceStatus, resource); drifted || err != nil {
						if err != nil {
							return nil, err
						}
						// resource drifted from state, Update
						description := fmt.Sprintf("Tool (%s/%s) drifted from the state, will be updated.", tool.Name, tool.InstanceID)
						changes = append(changes, generateUpdateAction(tool, description))
					} else {
						// resource is the same as the state, do nothing
						log.Debugf("Tool (%s/%s) is the same as the state, do nothing.", tool.Name, tool.InstanceID)
					}
				}
			}

			// delete the tool from the temporary state map since it's already been processed above
			tmpStatesMap.Delete(statemanager.GenerateStateKeyByToolNameAndInstanceID(tool.Name, tool.InstanceID))
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

	log.Debugf("Changes for the plan:")
	for i, c := range changes {
		log.Debugf("Change - %d/%d -> %s", i+1, len(changes), c.String())
	}

	return changes, nil
}

// GetChangesForDelete takes "State Manager" & "Config" then do some calculation and return a Plan to delete all plugins in the Config.
// All actions should be executed is included in this Plan.changes.
func GetChangesForDelete(smgr statemanager.Manager, cfg *configmanager.Config, isForceDelete bool) (changes []*Change, err error) {
	batchesOfTools, err := topologicalSort(cfg.Tools)
	if err != nil {
		return changes, err
	}

	log.Debug("isForce:", isForceDelete)
	for i := len(batchesOfTools) - 1; i >= 0; i-- {
		batch := batchesOfTools[i]
		for _, tool := range batch {
			if !isForceDelete {
				state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndInstanceID(tool.Name, tool.InstanceID))
				if state == nil {
					continue
				}
			}

			description := fmt.Sprintf("Tool (%s/%s) will be deleted.", tool.Name, tool.InstanceID)
			changes = append(changes, generateDeleteAction(tool, description))
		}
	}

	log.Debugf("Changes for the plan:")
	for i, c := range changes {
		log.Debugf("Change - %d/%d -> %s", i+1, len(changes), c.String())
	}

	return changes, nil
}

func GetChangesForDestroy(smgr statemanager.Manager, isForceDestroy bool) (changes []*Change, err error) {
	// rebuilding tools from config
	// because destroy will only be used when the config file is missing
	var tools configmanager.Tools
	for _, state := range smgr.GetStatesMap().ToList() {
		tool := &configmanager.Tool{
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
			if !isForceDestroy {
				state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndInstanceID(tool.Name, tool.InstanceID))
				if state == nil {
					continue
				}
			}
			description := fmt.Sprintf("Tool (%s/%s) will be deleted.", tool.Name, tool.InstanceID)
			changes = append(changes, generateDeleteAction(tool, description))
		}
	}

	log.Debugf("Changes for the plan:")
	for i, c := range changes {
		log.Debugf("Change - %d/%d -> %s", i+1, len(changes), c.String())
	}

	return changes, nil
}

// topologicalSortChangesInBatch returns a list of batches of changes, sorted by dependency defined in the config.
func topologicalSortChangesInBatch(changes []*Change) ([][]*Change, error) {

	// get tools from changes
	tools := getToolsFromChanges(changes)
	// map key to changes
	changesKeyMap := getChangesKeyMap(changes)

	batchesOfChanges := make([][]*Change, 0)

	// topological sort the tools based on dependency
	// tool in each batch does not have dependency on each other
	// note: maybe tools gotten from changes do not contain all the tools in the config.Tools
	// but it's still ok to call topologicalSort
	batchesOfTools, err := topologicalSort(tools)
	if err != nil {
		return batchesOfChanges, err
	}

	// group changes into batches based on the batches of tools
	// so that the changes in each batch do not have dependency on each other
	for _, batch := range batchesOfTools {
		changesOneBatch := make([]*Change, 0)
		// for each tool in the batch, find the change that matches it
		for _, tool := range batch {
			// only add the change that has the tool match with it
			if change, ok := changesKeyMap[tool.KeyWithNameAndInstanceID()]; ok {
				changesOneBatch = append(changesOneBatch, change)
			}
		}

		if len(changesOneBatch) > 0 {
			batchesOfChanges = append(batchesOfChanges, changesOneBatch)
		}
	}

	return batchesOfChanges, nil
}

func getToolsFromChanges(changes []*Change) configmanager.Tools {
	// use slice instead of map to keep the order of tools
	tools := make(configmanager.Tools, 0)
	// use map to record the tool that has been added to the slice
	toolsKeyMap := make(map[string]struct{})

	// get tools from changes avoiding duplicated tools
	for _, change := range changes {
		if _, ok := toolsKeyMap[change.Tool.KeyWithNameAndInstanceID()]; !ok {
			tools = append(tools, change.Tool)
			toolsKeyMap[change.Tool.KeyWithNameAndInstanceID()] = struct{}{}
		}
	}

	return tools
}

func getChangesKeyMap(changes []*Change) map[string]*Change {
	changesKeyMap := make(map[string]*Change)
	for _, change := range changes {
		changesKeyMap[change.Tool.KeyWithNameAndInstanceID()] = change
	}
	return changesKeyMap
}
