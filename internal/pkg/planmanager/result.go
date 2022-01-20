package planmanager

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/internal/pkg/util/log"
)

// HandleResult is used to Write the latest StatesMap to the Backend.
func (p *Plan) HandleResult(change *Change) error {
	if !change.Result.Succeeded {
		log.Errorf("Plugin %s %s failed.", change.Tool.Name, change.ActionName)
		return fmt.Errorf("plugin %s %s failed", change.Tool.Name, change.ActionName)
	}

	if change.ActionName == statemanager.ActionUninstall {
		p.smgr.DeleteState(getStateKeyFromTool(change.Tool))
		log.Successf("Plugin %s uninstall done.", change.Tool.Name)
		return p.smgr.Write(p.smgr.GetStatesMap().Format())
	}

	// install, reinstall
	var state = statemanager.NewState(
		change.Tool.Name,
		change.Tool.Plugin,
		[]string{},
		change.Tool.Options,
	)
	p.smgr.AddState(state)
	log.Successf("Plugin %s %s done.", change.Tool.Name, change.ActionName)
	return p.smgr.Write(p.smgr.GetStatesMap().Format())
}
