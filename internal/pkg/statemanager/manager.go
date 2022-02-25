package statemanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/pkg/util/log"
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

	GetState(key string) *State
	AddState(key string, state State) error
	UpdateState(key string, state State) error
	DeleteState(key string) error
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
	b, err := backend.GetBackend(backend.BackendLocal)
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

	tmpMap := make(map[string]State)
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

func (m *manager) GetState(key string) *State {
	m.statesMap.Load(key)
	if s, exist := m.statesMap.Load(key); exist {
		state, _ := s.(State)
		return &state
	}
	return nil
}

func (m *manager) AddState(key string, state State) error {
	m.statesMap.Store(key, state)
	return m.Write(m.GetStatesMap().Format())
}

func (m *manager) UpdateState(key string, state State) error {
	m.statesMap.Store(key, state)
	return m.Write(m.GetStatesMap().Format())
}

func (m *manager) DeleteState(key string) error {
	m.statesMap.Delete(key)
	return m.Write(m.GetStatesMap().Format())
}
