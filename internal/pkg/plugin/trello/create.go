package trello

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates Tello board and lists(todo/doing/done).
func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			setDefault,
			validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			createBoard,
			addTrelloSecret,
			cifile.PushCIFiles,
		},
		GetStatusOperation: getState,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(configmanager.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
