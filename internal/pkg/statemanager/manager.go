package statemanager

import (
	"fmt"
	"sync"

	"github.com/merico-dev/stream/internal/pkg/backend"
)

// Manager knows how to manage the States.
type Manager interface {
	backend.Backend
	// GetStates is used for generate data to Write to the backend.
	GetStates() States
	// SetStates is used for initialize States by the data Read from the backend.
	SetStates(states States)

	GetState(name string) *State
	AddState(state *State) error
	UpdateState(state *State)
	DeleteState(name string)
}

// manager is the default implement with Manager
type manager struct {
	mu sync.Mutex
	backend.Backend
	states States
}

func NewManager(backend backend.Backend) Manager {
	return &manager{
		Backend: backend,
		states:  make(States),
	}
}

func (m *manager) GetStates() States {
	return m.states
}

func (m *manager) SetStates(states States) {
	m.states = states
}

func (m *manager) GetState(name string) *State {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, exist := m.states[name]; exist {
		return s
	}
	return nil
}

func (m *manager) AddState(state *State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exist := m.states[state.Name]; exist {
		return fmt.Errorf("statemanager already exists")
	}
	m.states[state.Name] = state
	return nil
}

func (m *manager) UpdateState(state *State) {
	m.mu.Lock()
	m.states[state.Name] = state
	m.mu.Unlock()
}

func (m *manager) DeleteState(name string) {
	m.mu.Lock()
	delete(m.states, name)
	m.mu.Unlock()
}
