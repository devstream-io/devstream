package planmanager

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginengine"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// handleResult is used to Write the latest StatesMap to Backend.
func (p *Plan) handleResult(change *Change) error {
	if change.ActionName == statemanager.ActionUninstall {
		p.smgr.DeleteState(change.Tool.Name)
		return p.smgr.Write(p.smgr.GetStatesMap().Format())
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
	return p.smgr.Write(p.smgr.GetStatesMap().Format())
}

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func (p *Plan) generatePlanAccordingToConfig(statesMap *statemanager.StatesMap, cfg *configloader.Config) {
	for _, tool := range cfg.Tools {
		state := p.smgr.GetState(tool.Name)
		if state == nil {
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
				Action:     pluginengine.Install,
			})
			log.Printf("added a change: %s -> %s", tool.Name, statemanager.ActionInstall)
			continue
		}

		switch state.Status {
		case statemanager.StatusUninstalled:
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
				Action:     pluginengine.Install,
			})
			log.Printf("added a change: %s -> %s", tool.Name, statemanager.ActionInstall)
		case statemanager.StatusFailed:
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionReinstall,
				Action:     pluginengine.Reinstall,
			})
			log.Printf("added a change: %s -> %s", tool.Name, statemanager.ActionReinstall)
		}
		statesMap.Delete(tool.Name)
	}
}

// Some tools have already been installed, but they are no longer needed, so they need to be uninstalled
func (p *Plan) removeNoLongerNeededToolsFromPlan(statesMap *statemanager.StatesMap) {
	statesMap.Range(func(key, value interface{}) bool {
		p.Changes = append(p.Changes, &Change{
			Tool: &configloader.Tool{
				Name:    key.(string),
				Version: value.(*statemanager.State).Version,
				Options: value.(*statemanager.State).LastOperation.Metadata,
			},
			ActionName: statemanager.ActionUninstall,
			Action:     pluginengine.Uninstall,
		})
		log.Printf("added a change: %s -> %s", key.(string), statemanager.ActionUninstall)
		return true
	})
}
