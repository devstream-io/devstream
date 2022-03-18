package pluginengine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

// TestGetOutputs normal test, get Outputs by Statekey
func TestGetOutputs(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs":    map[string]interface{}{"boardId": "123", "todoListId": "456"}})

	defer func() {
		if err := smgr.DeleteState(key); err != nil {
			log.Errorf("failed to delete state %s.", key)
		}
	}()

	outputs, err := smgr.GetOutputs(key)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"boardId": "123", "todoListId": "456"}, outputs)
}

// TestGetOutputsIfEmpty exeption test, get Outputs by Statekey, but Outputs is empty
func TestGetOutputsIfEmpty(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456"})

	defer func() {
		if err := smgr.DeleteState(key); err != nil {
			log.Errorf("failed to delete state %s.", key)
		}
	}()

	outputs, err := smgr.GetOutputs(key)
	assert.EqualError(t, err, "cannot find outputs from state: default")
	assert.Equal(t, nil, outputs)
}

// TestGetOutputs exeption test, get Outputs by Statekey, but Statekey is wrong
func TestGetOutputsIfWrongKey(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)

	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs":    map[string]interface{}{"boardId": "123", "todoListId": "456"}})

	defer func() {
		if err := smgr.DeleteState(key); err != nil {
			log.Errorf("failed to delete state %s.", key)
		}
	}()

	outputs, err := smgr.GetOutputs("wrong_key")
	assert.EqualError(t, err, "cannot find state by key: wrong_key")
	assert.Equal(t, nil, outputs)
}

func initData(t *testing.T, smgr statemanager.Manager, resource map[string]interface{}) (statemanager.StateKey, statemanager.State) {
	keyReferred := statemanager.StateKey("default_trello")
	stateReferred := statemanager.State{
		Name:     "default",
		Plugin:   configloader.Plugin{Kind: "trello", Version: "0.2.0"},
		Options:  map[string]interface{}{"owner": "lfbdev", "repo": "golang-demo", "kanbanBoardName": "kanban-name"},
		Resource: resource,
	}
	err := smgr.AddState(keyReferred, stateReferred)
	assert.NoError(t, err)
	return keyReferred, stateReferred
}
