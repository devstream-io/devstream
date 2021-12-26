package planmanager

import (
	"log"
	"time"

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
			log.Printf("devstream.statesMap format error")
			return &Plan{Changes: make([]*Change, 0)}
		}
		for k, v := range tmpMap {
			statesMap.Store(k, v)
		}
		smgr.SetStatesMap(statesMap)
		log.Println("succeeded to initialize StatesMap")
	} else {
		log.Printf("failed to initialize StatesMap. Error: (%s). try to initialize the StatesMap", err)
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

// Execute will execute all changes included in the Plan and record results.
// All errors will be return.
func (p *Plan) Execute() []error {
	errors := make([]error, 0)
	log.Printf("changes count: %d", len(p.Changes))
	for i, c := range p.Changes {
		log.Printf("processing progress: %d/%d", i+1, len(p.Changes))
		log.Printf("processing: %s -> %s", c.Tool.Name, c.ActionName)
		// We will consider how to execute Action concurrently later.
		// It involves dependency management.
		succeeded, err := c.Action(c.Tool)
		if err != nil {
			errors = append(errors, err)
		}

		c.Result = &ChangeResult{
			Succeeded: succeeded,
			Error:     err,
			Time:      time.Now().Format(time.RFC3339),
		}

		err = p.handleResult(c)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
