package planmanager

import (
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// ActionFunc is a function that Do Action with a plugin. like:
// engine.Install() / engine.Reinstall() / engine.Uninstall()
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
