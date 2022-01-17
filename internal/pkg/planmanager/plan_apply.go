package planmanager

import (
	"log"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// Plan is an "Actions" plan, it includes all changes should be take with plugins.
type Plan struct {
	Changes []*Change
	smgr    statemanager.Manager
}

// NewPlan takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be execute is included in this Plan.changes.
func NewPlan(smgr statemanager.Manager, cfg *configloader.Config) *Plan {
	if cfg == nil {
		return &Plan{Changes: make([]*Change, 0)}
	}

	data, err := smgr.Read()
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

		log.Printf("Succeeded initializing StatesMap.")
	} else {
		log.Printf("Failed to initialize StatesMap. Error: (%s). Try to initialize the StatesMap.", err)
	}

	plan := &Plan{
		Changes: make([]*Change, 0),
		smgr:    smgr,
	}
	tmpStates := smgr.GetStatesMap().DeepCopy()
	plan.generatePlanAccordingToConfig(tmpStates, cfg)
	plan.removeNoLongerNeededToolsFromPlan(tmpStates)
	return plan
}
