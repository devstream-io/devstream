package devlake

import "github.com/devstream-io/devstream/internal/pkg/plugininstaller"

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		GetStateOperation: getDynamicState,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	return status, nil
}
