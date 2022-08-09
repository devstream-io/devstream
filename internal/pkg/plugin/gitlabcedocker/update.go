package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. create config and pre-handle operations
	opts, err := validateAndDefault(options)
	if err != nil {
		return nil, err
	}

	// reserve data when updated
	opts.RmDataAfterDelete = types.Bool(false)

	// 2. config install operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			dockerInstaller.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			dockerInstaller.DeleteAll,
			dockerInstaller.Install,
			showHelpMsg,
		},
		GetStateOperation: dockerInstaller.GetRunningState,
	}

	// 3. update and get status
	status, err := runner.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)

	return status, nil
}
