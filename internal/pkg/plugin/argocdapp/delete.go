package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. config install operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			kubectl.ProcessByContent("delete", "", argoCDAppTemplate),
		},
	}

	// 2. execute installer get status and error
	_, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}
