package devlakeconfig

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			validate,
			RenderAuthConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			DeleteConfig,
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	return true, nil
}
