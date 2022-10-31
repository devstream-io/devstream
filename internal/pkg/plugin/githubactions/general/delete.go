package general

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			setDefault,
			validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			//TODO(jiafeng meng): delete github secret
			cifile.DeleteCIFiles,
		},
		GetStatusOperation: getState,
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(configmanager.RawOptions(options))
	if err != nil {
		return false, err
	}

	return true, nil
}
