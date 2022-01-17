package planmanager

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// HandleResult is used to Write the latest StatesMap to the Backend.
func (p *Plan) HandleResult(change *Change) error {
	if !change.Result.Succeeded {
		log.Printf("Plugin %s %s failed.", change.Tool.Name, change.ActionName)
		return fmt.Errorf("plugin %s %s failed", change.Tool.Name, change.ActionName)
	}

	if change.ActionName == statemanager.ActionUninstall {
		p.smgr.DeleteState(getStateKeyFromTool(change.Tool))
		log.Printf("Plugin %s uninstall done.", change.Tool.Name)
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
	log.Printf("Plugin %s %s done.", change.Tool.Name, change.ActionName)
	return p.smgr.Write(p.smgr.GetStatesMap().Format())
}
