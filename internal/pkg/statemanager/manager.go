package statemanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/backend"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type ComponentAction string

const (
	ActionCreate ComponentAction = "Create"
	ActionUpdate ComponentAction = "Update"
	ActionDelete ComponentAction = "Delete"
)

// Manager knows how to manage the StatesMap.
type Manager interface {
	backend.Backend

	GetStatesMap() StatesMap

	GetState(key StateKey) *State
	AddState(key StateKey, state State) error
	UpdateState(key StateKey, state State) error
	DeleteState(key StateKey) error
	GetOutputs(key StateKey) (ResourceOutputs, error)
}

// manager is the default implement with Manager
type manager struct {
	backend.Backend
	statesMap StatesMap
}

var m *manager

// NewManager returns a new Manager and reads states through backend defined in config.
func NewManager(stateConfig configmanager.State) (Manager, error) {
	if m != nil {
		return m, nil
	}

	log.Debugf("The global manager m is not initialized.")

	// Get the backend from config
	// backend config has been validated in configmanager.Manager.LoadConfig()
	b, err := backend.GetBackend(stateConfig)
	if err != nil {
		log.Errorf("Failed to get the Backend: %s.", err)
		return nil, err
	}

	m = &manager{
		Backend:   b,
		statesMap: NewStatesMap(),
	}

	// Read the initial states data from backend
	data, err := b.Read()
	if err != nil {
		log.Debugf("Failed to read data from backend: %s.", err)
		return nil, err
	}

	tmpMap := make(map[StateKey]State)
	if err = yaml.Unmarshal(data, tmpMap); err != nil {
		log.Errorf("Failed to unmarshal the state of type < %s >, . error: %s.", stateConfig.Backend, err)
		log.Errorf("Reading the state failed, it might have been compromised/modified by someone other than DTM.")
		log.Errorf("The state is managed by DTM automatically. Please do not modify it yourself.")
		return nil, fmt.Errorf("state format error")
	}
	for k, v := range tmpMap {
		log.Debugf("Got a state from the backend: %s -> %v.", k, v)
		m.statesMap.Store(k, v)
	}

	return m, nil
}

func (m *manager) GetStatesMap() StatesMap {
	return m.statesMap
}

func (m *manager) GetState(key StateKey) *State {
	if s, exist := m.statesMap.Load(key); exist {
		state, _ := s.(State)
		return &state
	}
	return nil
}

// AddState adds a new state to the manager.
// If the state already exists, update it.
func (m *manager) AddState(key StateKey, state State) error {
	m.statesMap.Store(key, state)
	return m.Backend.Write(m.GetStatesMap().Format())
}

// UpdateState adds a new state to the manager.
// If the state already exists, update it.
// note: maybe it is duplicated with AddState
func (m *manager) UpdateState(key StateKey, state State) error {
	m.statesMap.Store(key, state)
	return m.Backend.Write(m.GetStatesMap().Format())
}

// DeleteState deletes a state from the manager.
// If the state does not exist, do nothing.
func (m *manager) DeleteState(key StateKey) error {
	m.statesMap.Delete(key)
	return m.Backend.Write(m.GetStatesMap().Format())
}

// GetOutputs is used to get the origin outputs of a toolName_InstanceID
func (m *manager) GetOutputs(key StateKey) (ResourceOutputs, error) {
	state := m.GetState(key)
	if state == nil {
		return nil, fmt.Errorf(`key (%s) not in state, it may be failed when "Create"`, key)
	}

	if value := state.ResourceStatus.GetOutputs(); value != nil {
		return value, nil
	}

	return nil, fmt.Errorf("outputs not in state %s", state.Name)
}
