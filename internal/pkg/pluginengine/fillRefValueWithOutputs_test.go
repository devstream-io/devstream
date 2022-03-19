package pluginengine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

// TestOutputsInState normal test, if the ref value is correct
func TestOutputsInState(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs":    map[string]interface{}{"boardId": "123", "todoListId": "456"}})

	defer func() {
		clearState(smgr, key)
	}()

	options := map[string]interface{}{
		"owner":      "lfbdev",
		"repo":       "golang-demo",
		"branch":     "main",
		"boardId":    "${{ default.trello.outputs.boardId }}",
		"todoListId": "${{ default.trello.outputs.todoListId }}"}
	expectResult := map[string]interface{}{"owner": "lfbdev",
		"repo":       "golang-demo",
		"branch":     "main",
		"boardId":    "123",
		"todoListId": "456"}

	err = fillRefValueWithOutputs(smgr, options)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectResult, options)
}

// TestRefInDeeperLayerState normal test, if when the ref key is in deeper layer
func TestRefInDeeperLayerState(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs": map[string]interface{}{
			"boardId":    "123",
			"todoListId": "456"}})

	defer func() {
		clearState(smgr, key)
	}()

	options := map[string]interface{}{
		"owner": "lfbdev",
		"repo":  "golang-demo", "branch": "main",
		"app": map[string]interface{}{
			"boardId":    "${{ default.trello.outputs.boardId }}",
			"todoListId": "${{ default.trello.outputs.todoListId }}"}}

	expectResult := map[string]interface{}{
		"owner":  "lfbdev",
		"repo":   "golang-demo",
		"branch": "main",
		"app": map[string]interface{}{
			"boardId":    "123",
			"todoListId": "456"}}

	err = fillRefValueWithOutputs(smgr, options)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectResult, options)
}

// TestOutputsEmpty exception test, if there is no outputs in state.
func TestOutputsEmpty(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{"boardId": "123", "todoListId": "456"})

	defer func() {
		clearState(smgr, key)
	}()

	options := map[string]interface{}{
		"owner":      "lfbdev",
		"repo":       "golang-demo",
		"branch":     "main",
		"boardId":    "${{ default.trello.outputs.boardId }}",
		"todoListId": "${{ default.trello.outputs.todoListId }}"}

	err = fillRefValueWithOutputs(smgr, options)

	assert.Equal(t, "cannot find outputs from state: default", err.Error())
}

// TestRefNotInOutput exception test, if the outputs do not contain ref
func TestRefNotInOutput(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs": map[string]interface{}{
			"boardId":    "123",
			"todoListId": "456"}})

	defer func() {
		clearState(smgr, key)
	}()

	options := map[string]interface{}{
		"owner":   "lfbdev",
		"repo":    "golang-demo",
		"branch":  "main",
		"boardId": "${{ default.trello.outputs.doesNotExistName }}"}

	err = fillRefValueWithOutputs(smgr, options)

	assert.Equal(t, "can not find doesNotExistName in dependency outputs", err.Error())
}

// TestRefLength exception test, if the ref length is correct.
func TestRefLength(t *testing.T) {
	smgr, err := statemanager.NewManager()
	assert.NoError(t, err)
	key, _ := initData(t, smgr, map[string]interface{}{
		"boardId":    "123",
		"todoListId": "456",
		"outputs": map[string]interface{}{
			"boardId":    "123",
			"todoListId": "456"}})

	defer func() {
		clearState(smgr, key)
	}()

	options := map[string]interface{}{
		"owner":   "lfbdev",
		"repo":    "golang-demo",
		"branch":  "main",
		"boardId": "${{ default.trello.outputs }}"}

	err = fillRefValueWithOutputs(smgr, options)

	assert.Equal(t, "incorrect output reference: default.trello.outputs", err.Error())
}

// clearState clear states when test ends
func clearState(smgr statemanager.Manager, keyReferred statemanager.StateKey) {
	if err := smgr.DeleteState(keyReferred); err != nil {
		log.Errorf("failed to delete state %s.", keyReferred)
	}
}
