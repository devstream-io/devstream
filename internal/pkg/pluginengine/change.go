package pluginengine

import (
	"fmt"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Change is a wrapper for a single Tool and its Action to be executed.
type Change struct {
	Tool        *configmanager.Tool
	ActionName  statemanager.ComponentAction
	Result      *ChangeResult
	Description string
}

// ChangeResult holds the result of a change action.
type ChangeResult struct {
	Succeeded   bool
	Error       error
	Time        string
	ReturnValue statemanager.ResourceStatus
}

func (c *Change) String() string {
	return fmt.Sprintf("\n{\n  ActionName: %s,\n  Tool: {Name: %s, InstanceID: %s}}\n}",
		c.ActionName, c.Tool.Name, c.Tool.InstanceID)
}

// execute changes the plan in batch.
// If any error occurs, it will stop executing the next batches and return the error.
func execute(smgr statemanager.Manager, changes []*Change, reverse bool) map[string]error {
	errorsMap := make(map[string]error)

	log.Info("Start executing the plan.")
	numOfChanges := len(changes)
	log.Infof("Changes count: %d.", numOfChanges)

	// get changes in batch
	// the changes in each batch do not have dependency on each other
	// but the changes from next batch have dependency on the changes from previous batch
	batchesOfChanges, err := topologicalSortChangesInBatch(changes)

	// for delete/destroy, the orders need to be reversed
	// so that the dependencies are deleted at last
	if reverse {
		for i, j := 0, len(batchesOfChanges)-1; i < j; i, j = i+1, j-1 {
			batchesOfChanges[i], batchesOfChanges[j] = batchesOfChanges[j], batchesOfChanges[i]
		}
	}

	if err != nil {
		log.Errorf("Failed to sort changes in batch: %s", err)
		errorsMap["dependency-analysis"] = err
		return errorsMap
	}

	currentChangeNum := 0
	for _, batch := range batchesOfChanges {
		for _, c := range batch {
			currentChangeNum += 1
			log.Separatorf("Processing progress: %d/%d.", currentChangeNum, numOfChanges)
			log.Infof("Processing: (%s/%s) -> %s ...", c.Tool.Name, c.Tool.InstanceID, c.ActionName)

			var succeeded bool
			var err error
			var returnValue statemanager.ResourceStatus

			log.Debugf("Tool's raw changes are: %v.", c.Tool.Options)

			errs := HandleOutputsReferences(smgr, c.Tool.Options)
			if len(errs) != 0 {
				succeeded = false

				for _, e := range errs {
					log.Errorf("Error: %s.", e)
				}
				log.Errorf("The outputs reference in tool (%s/%s) can't be resolved. Please double check your config.", c.Tool.Name, c.Tool.InstanceID)

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
				key := fmt.Sprintf("%s/%s/%s", c.Tool.Name, c.Tool.InstanceID, c.ActionName)
				log.Errorf("%s/%s %s failed with error: %s", c.Tool.Name, c.Tool.InstanceID, c.ActionName, err)
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

		// abort next batches if any error occurred in this batch
		if len(errorsMap) != 0 {
			break
		}
	}

	if len(errorsMap) != 0 {
		log.Separatorf("Processing aborted.")
	} else {
		log.Separatorf("Processing done.")
	}

	return errorsMap
}

// handleResult is used to Write the latest StatesMap to the Backend.
func handleResult(smgr statemanager.Manager, change *Change) error {
	log.Debugf("Start -> StatesMap now is:\n%s", string(smgr.GetStatesMap().Format()))
	defer func() {
		log.Debugf("End -> StatesMap now is:\n%s", string(smgr.GetStatesMap().Format()))
	}()

	if !change.Result.Succeeded {
		// do nothing when the change failed
		return nil
	}

	if change.ActionName == statemanager.ActionDelete {
		key := statemanager.GenerateStateKeyByToolNameAndInstanceID(change.Tool.Name, change.Tool.InstanceID)
		log.Infof("Prepare to delete '%s' from States.", key)
		err := smgr.DeleteState(key)
		if err != nil {
			log.Debugf("Failed to delete state %s: %s.", key, err)
			return err
		}
		log.Successf("Tool (%s/%s) delete done.", change.Tool.Name, change.Tool.InstanceID)
		return nil
	}

	key := statemanager.GenerateStateKeyByToolNameAndInstanceID(change.Tool.Name, change.Tool.InstanceID)
	state := statemanager.State{
		Name:           change.Tool.Name,
		InstanceID:     change.Tool.InstanceID,
		DependsOn:      change.Tool.DependsOn,
		Options:        change.Tool.Options,
		ResourceStatus: change.Result.ReturnValue,
	}
	err := smgr.AddState(key, state)
	if err != nil {
		log.Debugf("Failed to add state %s: %s.", key, err)
		return err
	}
	log.Successf("Tool (%s/%s) %s done.", change.Tool.Name, change.Tool.InstanceID, change.ActionName)
	return nil
}
