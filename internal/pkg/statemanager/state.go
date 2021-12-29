package statemanager

import (
	"bytes"
	"log"

	"gopkg.in/yaml.v3"
)

type ComponentStatus string
type ComponentAction string

const (
	// We should delete the state of the "uninstalled" tool at States.
	StatusUninstalled ComponentStatus = "uninstalled"
	// We use StatusInstalled when a plugin is installed but we don't know its status is "running" or "failed".
	// For example: We try to uninstall a plugin but failed for some reason.
	StatusInstalled ComponentStatus = "installed"
	StatusRunning   ComponentStatus = "running"
	StatusFailed    ComponentStatus = "failed"
)

const (
	ActionInstall   ComponentAction = "install"
	ActionReinstall ComponentAction = "reinstall"
	ActionUninstall ComponentAction = "uninstall"
)

type States map[string]*State

// State is the single component's state.
type State struct {
	Name          string          `yaml:"name"`
	Version       string          `yaml:"version"`
	Dependencies  []string        `yaml:"dependencies"`
	Status        ComponentStatus `yaml:"status"`
	LastOperation *Operation      `yaml:"lastOperation"`
}

type Operation struct {
	Action   ComponentAction        `yaml:"action"`
	Time     string                 `yaml:"time"`
	Metadata map[string]interface{} `yaml:"metadata"`
}

func NewState(name, version string, deps []string, status ComponentStatus, lastOpr *Operation) *State {
	return &State{
		Name:          name,
		Version:       version,
		Dependencies:  deps,
		Status:        status,
		LastOperation: lastOpr,
	}
}

func (s States) DeepCopy() States {
	if s == nil {
		return nil
	}
	newStates := make(States)
	for k, v := range s {
		newStates[k] = v
	}
	return newStates
}

func (s States) Format() []byte {
	if len(s) == 0 {
		return []byte{}
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(&s)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}
