package statemanager

import (
	"os"
	"reflect"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/backend/local"
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
	lastOperation := &Operation{
		Action: ActionInstall,
		Time:   time.Now().Format("2006-01-02_15:04:05"),
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}
	return NewState("argocd", "v0.0.1", nil, StatusRunning, lastOperation)
}

func TestManager_State(t *testing.T) {
	setup(t)
	stateA := newState()
	err := smgr.AddState(stateA)
	if err != nil {
		t.Fatal(err)
	}

	stateB := smgr.GetState("argocd")

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
	_ = smgr.AddState(stateA)
	if err := smgr.Write(smgr.GetStates().Format()); err != nil {
		t.Error("Failed to Write States to disk")
	}
}

func TestManager_Read(t *testing.T) {
	TestManager_Write(t)
	data, err := smgr.Read()
	if err != nil {
		t.Error(err)
	}

	var oldSs = make(States)
	if err := yaml.Unmarshal(data, oldSs); err != nil {
		t.Error(err)
	}

	smgr.SetStates(oldSs)
	newSs := smgr.GetStates()
	if !reflect.DeepEqual(smgr.GetStates(), oldSs) {
		t.Errorf("expect old States == new States, but got oldSs: %v and newSs: %v", oldSs, newSs)
	}

	teardown()
}

func teardown() {
	os.Remove(local.DefaultStateFile)
}
