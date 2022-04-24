package statemanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/backend"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
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
	GetOutputs(key StateKey) (interface{}, error)
}

// manager is the default implement with Manager
type manager struct {
	backend.Backend
	statesMap StatesMap
}

var m *manager

func NewManager() (Manager, error) {
	if m != nil {
		return m, nil
	}

	log.Debugf("The global manager m is not initialized.")

	// use default local backend for now.
	b, err := backend.GetBackend(backend.Local)
	if err != nil {
		log.Errorf("Failed to get the Backend: %s.", err)
		return nil, err
	}

	m = &manager{
		Backend:   b,
		statesMap: NewStatesMap(),
	}

	// Read the initial states data
	data, err := b.Read()
	if err != nil {
		log.Debugf("Failed to read data from backend: %s.", err)
		return nil, err
	}

	tmpMap := make(map[StateKey]State)
	if err = yaml.Unmarshal(data, tmpMap); err != nil {
		log.Errorf("Failed to unmarshal the state file < %s >. error: %s.", local.DefaultStateFile, err)
		log.Errorf("Reading the state file failed, it might have been compromised/modified by someone other than DTM.")
		log.Errorf("The state file is managed by DTM automatically. Please do not modify it yourself.")
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
	m.statesMap.Load(key)
	if s, exist := m.statesMap.Load(key); exist {
		state, _ := s.(State)
		return &state
	}
	return nil
}

func (m *manager) AddState(key StateKey, state State) error {
	m.statesMap.Store(key, state)
	return m.Write(m.GetStatesMap().Format())
}

func (m *manager) UpdateState(key StateKey, state State) error {
	m.statesMap.Store(key, state)
	return m.Write(m.GetStatesMap().Format())
}

func (m *manager) DeleteState(key StateKey) error {
	m.statesMap.Delete(key)
	return m.Write(m.GetStatesMap().Format())
}

func (m *manager) GetOutputs(key StateKey) (interface{}, error) {
	state := m.GetState(key)
	if state == nil {
		return nil, fmt.Errorf("key (%s) not in state, it may be failed when \"Create\"", key)
	}

	if value, ok := state.Resource["outputs"]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("outputs not in state %s", state.Name)
}
