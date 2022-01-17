package statemanager

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
	"github.com/merico-dev/stream/internal/pkg/configloader"
)

var smgr Manager

// setup is used to initialize smgr.
func setup(t *testing.T) {
	b, err := backend.GetBackend("local")
	if err != nil {
		t.Fatal("failed to get backend.")
	}

	smgr = NewManager(b)
}

func newState() *State {
	return NewState("argocd-prod", configloader.Plugin{Kind: "argocd", Version: "v0.0.1"}, nil, nil)
}

func TestManager_State(t *testing.T) {
	setup(t)
	stateA := newState()
	smgr.AddState(stateA)

	stateB := smgr.GetState("argocd-prod_argocd")

	if !reflect.DeepEqual(stateA, stateB) {
		t.Errorf("expect stateB == stateA, but got stateA: %v and stateB: %v", stateA, stateB)
	}

	smgr.DeleteState("argocd")
	if smgr.GetState("argocd") != nil {
		t.Error("DeleteState failed")
	}
}

func TestManager_Write(t *testing.T) {
	setup(t)
	stateA := newState()
	smgr.AddState(stateA)
	if err := smgr.Write(smgr.GetStatesMap().Format()); err != nil {
		t.Error("Failed to Write StatesMap to disk")
	}
}

func TestManager_Read(t *testing.T) {
	TestManager_Write(t)
	data, err := smgr.Read()
	if err != nil {
		t.Error(err)
	}

	var oldSs = NewStatesMap()
	if err := yaml.Unmarshal(data, oldSs); err != nil {
		t.Error(err)
	}

	smgr.SetStatesMap(oldSs)
	newSs := smgr.GetStatesMap()
	if !reflect.DeepEqual(smgr.GetStatesMap(), oldSs) {
		t.Errorf("expect old StatesMap == new StatesMap, but got oldSs: %v and newSs: %v", oldSs, newSs)
	}

	teardown()
}

func teardown() {
	_ = os.Remove(local.DefaultStateFile)
}
