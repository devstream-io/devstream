package kubeprometheus

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			helm.SetDefaultConfig(&defaultHelmConfig),
			helm.Validate,
		},
		ExecuteOperations: helm.DefaultDeleteOperations,
	}
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
