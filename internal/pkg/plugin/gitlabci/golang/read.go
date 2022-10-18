package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			ci.SetDefaultConfig(gitlabci.DefaultCIOptions),
			setCIContent,
			ci.Validate,
		},
		GetStatusOperation: ci.GetCIFileStatus,
	}

	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}

	log.Debugf("Return map: %v", status)
	return status, nil
}
