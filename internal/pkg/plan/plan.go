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
}

// ExecutionResult holds the execution result with a single plugin.
// The map's key is a plugin's name.
type ExecutionResult map[string]*Change

// Execute will execute all changes included in the Plan and put the execution result into a ExecutionResult channel.
func (p *Plan) Execute(execResultChan chan<- ExecutionResult) {
	for _, c := range p.Changes {
		// We will consider how to execute Action concurrently later.
		// It involves dependency management.
		succeeded, err := c.Action(c.Tool)
		c.Result = &ChangeResult{
			Succeeded: succeeded,
			Error:     err,
			Time:      time.Now().Format("2006-01-02_15:04:05"),
		}
		execResultChan <- ExecutionResult{c.Tool.Name: c}
	}
	close(execResultChan)
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

// MakePlan takes "State Manager" & "Config" then do some calculate and return a Plan.
// All actions should be execute is included in this Plan.changes.
func MakePlan(smgr statemanager.Manager, cfg *config.Config) *Plan {
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

	plan := &Plan{Changes: make([]*Change, 0)}
	tmpStates := smgr.GetStates().DeepCopy()

	generatePlanAccordingToConfig(smgr, plan, tmpStates, cfg)
	removeNoLongerNeededToolsFromPlan(plan, tmpStates)

	return plan
}

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func generatePlanAccordingToConfig(smgr statemanager.Manager, plan *Plan, states statemanager.States, cfg *config.Config) {
	for _, tool := range cfg.Tools {
		state := smgr.GetState(tool.Name)
		if state == nil {
			plan.Changes = append(plan.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     plugin.Install,
			})
			continue
		}

		switch state.Status {
		case statemanager.StatusUninstalled:
			plan.Changes = append(plan.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionInstall,
				Action:     plugin.Install,
			})
		case statemanager.StatusFailed:
			plan.Changes = append(plan.Changes, &Change{
				Tool:       &tool,
				ActionName: statemanager.ActionReinstall,
				Action:     plugin.Reinstall,
			})
		}
		delete(states, tool.Name)
	}
}

// Some tools have already been installed, but they are no longer needed, so they need to be uninstalled
func removeNoLongerNeededToolsFromPlan(plan *Plan, states statemanager.States) {
	for _, state := range states {
		plan.Changes = append(plan.Changes, &Change{
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
