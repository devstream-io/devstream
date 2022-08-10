package kubeprometheus

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create will create prometheus in k8s
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			helm.Validate,
		},
		ExecuteOperations:   helm.DefaultCreateOperations,
		TerminateOperations: helm.DefaultTerminateOperations,
		GetStateOperation:   helm.GetPluginAllState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
