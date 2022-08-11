package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			renderConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			dockerInstaller.ComposeDown,
		},
		GetStateOperation: dockerInstaller.ComposeState,
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}
