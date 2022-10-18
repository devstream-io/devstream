package generic

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			ci.SetDefaultConfig(gitlabci.DefaultCIOptions),
			ci.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			ci.PushCIFiles,
		},
		GetStatusOperation: ci.GetCIFileStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
