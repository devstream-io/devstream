package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/goclient"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			goclient.Validate,
		},
		GetStatusOperation: goclient.GetState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	return status, nil
}
