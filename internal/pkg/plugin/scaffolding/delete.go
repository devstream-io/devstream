package scaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/reposcaffolding"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. config delete operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			reposcaffolding.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			reposcaffolding.DeleteRepo,
		},
	}
	_, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
