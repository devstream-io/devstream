package tekton

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config read operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			helm.SetDefaultConfig(&defaultHelmConfig),
			helm.Validate,
		},
		GetStatusOperation: helm.GetPluginAllState,
	}

	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}

	log.Debugf("Return map: %v", status)
	return status, nil
}
