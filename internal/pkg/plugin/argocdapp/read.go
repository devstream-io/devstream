package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			validate,
		},
		GetStatusOperation: getDynamicState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	return status, nil
}
