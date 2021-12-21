package planmanager

import (
	"log"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginengine"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// Plan is an "Actions" plan, it includes all changes should be take with plugins.
type Plan struct {
	Changes []*Change
	smgr    statemanager.Manager
}

// ActionFunc is a function that Do Action with a plugin. like:
// plugin.Install() / plugin.Reinstall() / plugin.Uninstall()
type ActionFunc func(tool *configloader.Tool) (bool, error)

// Change is a wrapper with a single Tool and its Action should be execute.
type Change struct {
	Tool       *configloader.Tool
	ActionName statemanager.ComponentAction
	Action     ActionFunc
	Result     *ChangeResult
}

// ChangeResult holds the result with a change action.
type ChangeResult struct {
	Succeeded bool
	Error     error
	Time      string
}

// NewPlan takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be execute is included in this Plan.changes.
func NewPlan(smgr statemanager.Manager, cfg *configloader.Config) *Plan {
	if cfg == nil {
		return &Plan{Changes: make([]*Change, 0)}
	}

	data, err := smgr.Read()
	if err == nil {
		states := make(statemanager.States)
		if err := yaml.Unmarshal(data, states); err != nil {
			log.Printf("devstream.states format error")
			return &Plan{Changes: make([]*Change, 0)}
		}
		smgr.SetStates(states)
		log.Println("succeeded to initialize States")
	} else {
		log.Printf("failed to initialize States. Error: (%s). try to initialize the States", err)
	}

	plan := &Plan{
		Changes: make([]*Change, 0),
		smgr:    smgr,
	}
	tmpStates := smgr.GetStates().DeepCopy()

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
		log.Printf("procprocessing progress: %d/%d", i+1, len(p.Changes))
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

// handleResult is used to Write the latest States to Backend.
func (p *Plan) handleResult(change *Change) error {
	if change.ActionName == statemanager.ActionUninstall {
		p.smgr.DeleteState(change.Tool.Name)
		return p.smgr.Write(p.smgr.GetStates().Format())
	}

	var state = statemanager.NewState(
		change.Tool.Name,
		change.Tool.Version,
		[]string{},
		statemanager.StatusRunning,
		&statemanager.Operation{
			Action:   change.ActionName,
			Time:     change.Result.Time,
			Metadata: change.Tool.Options,
		},
	)

	if change.Result.Error != nil {
		state.Status = statemanager.StatusFailed
		log.Printf("=== plugin %s process failed ===", change.Tool.Name)
	} else {
		log.Printf("=== plugin %s process done ===", change.Tool.Name)
	}

	p.smgr.AddState(state)
	return p.smgr.Write(p.smgr.GetStates().Format())
}

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func (p *Plan) generatePlanAccordingToConfig(states statemanager.States, cfg *configloader.Config) {
	for _, tool := range cfg.Tools {
		state := p.smgr.GetState(tool.Name)
		if state == nil {
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     pluginengine.Install,
			})
			continue
		}

		switch state.Status {
		case statemanager.StatusUninstalled:
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     pluginengine.Install,
			})
		case statemanager.StatusFailed:
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionReinstall,
				Action:     pluginengine.Reinstall,
			})
		}
		delete(states, tool.Name)
	}
}

// Some tools have already been installed, but they are no longer needed, so they need to be uninstalled
func (p *Plan) removeNoLongerNeededToolsFromPlan(states statemanager.States) {
	for _, state := range states {
		p.Changes = append(p.Changes, &Change{
			Tool: &configloader.Tool{
				Name:    state.Name,
				Version: state.Version,
				Options: state.LastOperation.Metadata,
			},
			ActionName: statemanager.ActionUninstall,
			Action:     pluginengine.Uninstall,
		})
	}
}
