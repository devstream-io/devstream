package plan

import (
	"log"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/plugin"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// Plan is an "Actions" plan, it includes all changes should be take with plugins.
type Plan struct {
	Changes []*Change
	smgr    statemanager.Manager
}

// ActionFunc is a function that Do Action with a plugin. like:
// plugin.Install() / plugin.Reinstall() / plugin.Uninstall()
type ActionFunc func(tool *config.Tool) (bool, error)

// Change is a wrapper with a single Tool and its Action should be execute.
type Change struct {
	Tool       *config.Tool
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
func NewPlan(smgr statemanager.Manager, cfg *config.Config) *Plan {
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
	}
	log.Printf("failed to initialize States. %s", err)
	log.Println("try to initialize the States.")

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
	for _, c := range p.Changes {
		// We will consider how to execute Action concurrently later.
		// It involves dependency management.
		succeeded, err := c.Action(c.Tool)
		if err != nil {
			errors = append(errors, err)
		}
		c.Result.Succeeded = succeeded
		c.Result.Error = err
		c.Result.Time = time.Now().Format(time.RFC3339)

		err = p.handleResult(c)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// handleResult is used to Write the latest States to Backend.
func (p *Plan) handleResult(change *Change) error {
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
		log.Printf("=== plugin %s installation failed ===", change.Tool.Name)
	} else {
		log.Printf("=== plugin %s installation done ===", change.Tool.Name)
	}

	p.smgr.AddState(state)
	return p.smgr.Write(p.smgr.GetStates().Format())
}

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func (p *Plan) generatePlanAccordingToConfig(states statemanager.States, cfg *config.Config) {
	for _, tool := range cfg.Tools {
		state := p.smgr.GetState(tool.Name)
		if state == nil {
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     plugin.Install,
			})
			continue
		}

		switch state.Status {
		case statemanager.StatusUninstalled:
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     plugin.Install,
			})
		case statemanager.StatusFailed:
			p.Changes = append(p.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionReinstall,
				Action:     plugin.Reinstall,
			})
		}
		delete(states, tool.Name)
	}
}

// Some tools have already been installed, but they are no longer needed, so they need to be uninstalled
func (p *Plan) removeNoLongerNeededToolsFromPlan(states statemanager.States) {
	for _, state := range states {
		p.Changes = append(p.Changes, &Change{
			Tool: &config.Tool{
				Name:    state.Name,
				Version: state.Version,
				Options: state.LastOperation.Metadata,
			},
			ActionName: statemanager.ActionUninstall,
			Action:     plugin.Uninstall,
		})
	}
}
