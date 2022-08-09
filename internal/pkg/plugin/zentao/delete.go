package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/goclient"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. config install operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			goclient.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			goclient.Delete,
		},
	}

	// 2. execute installer get status and error
	_, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}
