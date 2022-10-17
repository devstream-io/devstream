package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates ArgoCD with provided options.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			helm.SetDefaultConfig(&defaultHelmConfig),
			helm.Validate,
		},
		ExecuteOperations:   helm.DefaultCreateOperations,
		TerminateOperations: helm.DefaultTerminateOperations,
		GetStatusOperation:  helm.GetAllResourcesStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
