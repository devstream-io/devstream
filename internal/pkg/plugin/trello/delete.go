package trello

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
)

// Create creates Tello board and lists(todo/doing/done).
func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			setDefault,
			validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			deleteBoard,
			cifile.DeleteCIFiles,
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(configmanager.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}
