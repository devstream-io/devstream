package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			renderConfig,
		},
		ExecuteOperations: installer.ExecuteOperations{
			dockerInstaller.ComposeDown,
		},
		GetStatusOperation: dockerInstaller.ComposeStatus,
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}
