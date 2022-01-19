package planmanager

import (
	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/internal/pkg/util/log"
)

// NewDeletePlan takes "State Manager" & "Config" then do some calculation and return a Plan to delete all plugins in the Config.
// All actions should be execute is included in this Plan.changes.
func NewDeletePlan(smgr statemanager.Manager, cfg *configloader.Config) *Plan {
	if cfg == nil {
		return &Plan{Changes: make([]*Change, 0)}
	}

	data, err := smgr.Read()
	// TODO(ironcore864): duplicated code; needs to be refactored.
	if err == nil {
		statesMap := statemanager.NewStatesMap()
		tmpMap := make(map[string]*statemanager.State)
		if err := yaml.Unmarshal(data, tmpMap); err != nil {
			log.Fatalf("Devstream.statesMap format error.")
			return &Plan{Changes: make([]*Change, 0)}
		}
		for k, v := range tmpMap {
			statesMap.Store(k, v)
		}
		smgr.SetStatesMap(statesMap)
		log.Success("Succeeded initializing StatesMap.")
	} else {
		log.Errorf("Failed to initialize StatesMap. Error: (%s). Try to initialize the StatesMap.", err)
	}

	plan := &Plan{
		Changes: make([]*Change, 0),
		smgr:    smgr,
	}
	tmpStates := smgr.GetStatesMap().DeepCopy()
	plan.generatePlanForDelete(tmpStates, cfg)
	return plan
}
