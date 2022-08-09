package openldap

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. config delete operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			helm.Validate,
		},
		ExecuteOperations: helm.DefaultDeleteOperations,
	}
	_, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
