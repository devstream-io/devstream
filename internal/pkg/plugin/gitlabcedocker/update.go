package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. create config and pre-handle operations
	opts, err := preHandleOptions(options)
	if err != nil {
		return nil, err
	}

	opts.setGitLabURL()

	// 2. config install operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			dockerInstaller.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			dockerInstaller.InstallOrUpdate,
			showGitLabURL,
		},
		GetStatusOperation: dockerInstaller.GetRunningState,
	}

	// 3. update and get status
	status, err := runner.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)

	return status, nil
}
