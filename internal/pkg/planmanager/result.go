package planmanager

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// HandleResult is used to Write the latest StatesMap to the Backend.
func (p *Plan) HandleResult(change *Change) error {
	// uninstall succeeded
	if change.ActionName == statemanager.ActionUninstall && change.Result.Succeeded {
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

	// uninstall failed
	if change.ActionName == statemanager.ActionUninstall && !change.Result.Succeeded {
		state.Status = statemanager.StatusInstalled
		log.Printf("Plugin %s uninstall failed.", change.Tool.Name)
	} else if !change.Result.Succeeded {
		// install or reinstall failed
		state.Status = statemanager.StatusFailed
		log.Printf("Plugin %s (re)install failed.", change.Tool.Name)
	} else {
		// install or reinstall succeeded
		state.Status = statemanager.StatusInstalled
		log.Printf("Plugin %s process done.", change.Tool.Name)
	}

	p.smgr.AddState(state)
	return p.smgr.Write(p.smgr.GetStatesMap().Format())
}
