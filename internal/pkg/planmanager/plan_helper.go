package planmanager

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func (p *Plan) generatePlanAccordingToConfig(statesMap *statemanager.StatesMap, cfg *configloader.Config) {
	for _, tool := range cfg.Tools {
		state := p.smgr.GetState(tool.Name)
		if state == nil {
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionInstall)
			continue
		}

		switch state.Status {
		case statemanager.StatusUninstalled:
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionInstall)
		case statemanager.StatusFailed:
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionReinstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionReinstall)
		case statemanager.StatusRunning:
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionInstall)
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
		})
		log.Printf("Change added: %s -> %s", key.(string), statemanager.ActionUninstall)
		return true
	})
}

// generatePlanForDelete is to create a plan that deletes all the Tools in cfg
func (p *Plan) generatePlanForDelete(statesMap *statemanager.StatesMap, cfg *configloader.Config) {
	// reverse loop, a hack to solve dependency issues when uninstalling
	for i := len(cfg.Tools) - 1; i >= 0; i-- {
		tool := cfg.Tools[i]
		state := p.smgr.GetState(tool.Name)
		if state == nil || state.Status != statemanager.StatusInstalled {
			continue
		}

		p.Changes = append(p.Changes, &Change{
			Tool:       tool.DeepCopy(),
			ActionName: statemanager.ActionUninstall,
		})
		log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionUninstall)
		statesMap.Delete(tool.Name)
	}
}
