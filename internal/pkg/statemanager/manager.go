package statemanager

import (
	"github.com/merico-dev/stream/internal/pkg/backend"
)

// Manager knows how to manage the StatesMap.
type Manager interface {
	backend.Backend
	// GetStatesMap is used for generate data to Write to the backend.
	GetStatesMap() *StatesMap
	// SetStatesMap is used for initialize StatesMap by the data Read from the backend.
	SetStatesMap(states *StatesMap)

	GetState(name string) *State
	AddState(state *State)
	UpdateState(state *State)
	DeleteState(name string)
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

func (m *manager) GetState(name string) *State {
	m.statesMap.Load(name)
	if s, exist := m.statesMap.Load(name); exist {
		return s.(*State)
	}
	return nil
}

func (m *manager) AddState(state *State) {
	m.statesMap.Store(state.Name, state)
}

func (m *manager) UpdateState(state *State) {
	m.statesMap.Store(state.Name, state)
}

func (m *manager) DeleteState(name string) {
	m.statesMap.Delete(name)
}
