package statemanager

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/backend"
)

// Manager knows how to manage the StatesMap.
type Manager interface {
	backend.Backend
	// GetStatesMap returns the state manager's map of states, used to generate data to Write to the backend.
	GetStatesMap() *StatesMap
	// SetStatesMap sets the state manager's map of states, used to initialize the StatesMap by the data Read from the backend.
	SetStatesMap(states *StatesMap)

	GetState(key string) *State
	AddState(state *State)
	UpdateState(state *State)
	DeleteState(key string)
}

// manager is the default implement with Manager
type manager struct {
	backend.Backend
	statesMap *StatesMap
}

func NewManager(backend backend.Backend) Manager {
	return &manager{
		Backend:   backend,
		statesMap: NewStatesMap(),
	}
}

func (m *manager) GetStatesMap() *StatesMap {
	return m.statesMap
}

func (m *manager) SetStatesMap(statesMap *StatesMap) {
	m.statesMap = statesMap
}

func (m *manager) GetState(key string) *State {
	m.statesMap.Load(key)
	if s, exist := m.statesMap.Load(key); exist {
		return s.(*State)
	}
	return nil
}

func (m *manager) AddState(state *State) {
	stateKey := getStateKeyFromState(state)
	m.statesMap.Store(stateKey, state)
}

func (m *manager) UpdateState(state *State) {
	stateKey := getStateKeyFromState(state)
	m.statesMap.Store(stateKey, state)
}

func (m *manager) DeleteState(key string) {
	m.statesMap.Delete(key)
}

func getStateKeyFromState(s *State) string {
	return fmt.Sprintf("%s_%s", s.Name, s.Plugin.Kind)
}
