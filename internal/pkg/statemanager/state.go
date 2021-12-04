package statemanager

import (
	"bytes"
	"log"

	"gopkg.in/yaml.v3"
)

type CompStatus string
type CompAction string

const (
	StatusUninstalled CompStatus = "uninstalled"
	StatusInstalled   CompStatus = "installed"
	StatusRunning     CompStatus = "running"
	StatusFailed      CompStatus = "Failed"
)

const (
	ActionInstall   CompAction = "install"
	ActionReinstall CompAction = "reinstall"
	ActionUninstall CompAction = "uninstall"
)

type States map[string]*State

// State is the single component's state.
type State struct {
	Name          string     `yaml:"name"`
	Version       string     `yaml:"version"`
	Dependencies  []string   `yaml:"dependencies"`
	Status        CompStatus `yaml:"status"`
	LastOperation *Operation `yaml:"lastOperation"`
}

type Operation struct {
	Action   CompAction             `yaml:"action"`
	Time     string                 `yaml:"time"`
	Metadata map[string]interface{} `yaml:"metadata"`
}

func NewState(name, version string, deps []string, status CompStatus, lastOpr *Operation) *State {
	return &State{
		Name:          name,
		Version:       version,
		Dependencies:  deps,
		Status:        status,
		LastOperation: lastOpr,
	}
}

func (s States) Format() []byte {
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
