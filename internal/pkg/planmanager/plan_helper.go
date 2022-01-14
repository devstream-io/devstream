package planmanager

import (
	"log"
	"reflect"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func drifted(t *configloader.Tool, s *statemanager.State) bool {
	return !reflect.DeepEqual(t.Options, s.Metadata) || !reflect.DeepEqual(t.Plugin, s.Plugin)
}

// generatePlanAccordingToConfig is to filter all the Tools in cfg that need some actions
func (p *Plan) generatePlanAccordingToConfig(statesMap *statemanager.StatesMap, cfg *configloader.Config) {
	for _, tool := range cfg.Tools {
		state := p.smgr.GetState(getStateKeyFromTool(&tool))
		if state == nil {
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionInstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionInstall)
			continue
		}

		if drifted(&tool, state) {
			p.Changes = append(p.Changes, &Change{
				Tool:       tool.DeepCopy(),
				ActionName: statemanager.ActionReinstall,
			})
			log.Printf("Change added: %s -> %s", tool.Name, statemanager.ActionReinstall)
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
				Plugin:  value.(*statemanager.State).Plugin,
				Options: value.(*statemanager.State).Metadata,
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
		state := p.smgr.GetState(getStateKeyFromTool(&tool))
		if state == nil {
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
