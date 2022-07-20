package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	ArgocdDefaultNamespace = "argocd"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config read operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			defaultMissedOption,
			helm.Validate,
		},
		GetStatusOperation: helm.GetPluginStaticStateWrapper(defaultDeploymentList),
	}

	// 2. get plugin status
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
