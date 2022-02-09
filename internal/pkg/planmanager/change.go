package planmanager

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

// Change is a wrapper with a single Tool and its Action should be execute.
type Change struct {
	Tool       *configloader.Tool
	ActionName statemanager.ComponentAction
	Result     *ChangeResult
}

// ChangeResult holds the result with a change action.
type ChangeResult struct {
	Succeeded bool
	Error     error
	Time      string
}

func (c *Change) String() string {
	return fmt.Sprintf("\n{\n  ActionName: %s,\n  Tool: {Name: %s, Plugin: {Kind: %s, Version: %s}}\n}",
		c.ActionName, c.Tool.Name, c.Tool.Plugin.Kind, c.Tool.Plugin.Version)
}

func getStateKeyFromTool(t *configloader.Tool) string {
	return fmt.Sprintf("%s_%s", t.Name, t.Plugin.Kind)
}
