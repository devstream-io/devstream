package statemanager

import "github.com/merico-dev/stream/internal/pkg/configloader"

type ComponentAction string

const (
	ActionInstall   ComponentAction = "Install"
	ActionReinstall ComponentAction = "Reinstall"
	ActionUninstall ComponentAction = "Uninstall"
)

// State is the single component's state.
type State struct {
	Name         string                 `yaml:"name"`
	Plugin       configloader.Plugin    `yaml:"plugin"`
	Dependencies []string               `yaml:"dependencies"`
	Metadata     map[string]interface{} `yaml:"metadata"`
}

func NewState(name string, plugin configloader.Plugin, deps []string, metadata map[string]interface{}) *State {
	return &State{
		Name:         name,
		Plugin:       plugin,
		Dependencies: deps,
		Metadata:     metadata,
	}
}
