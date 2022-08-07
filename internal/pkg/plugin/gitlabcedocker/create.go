package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
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
		TerminateOperations: []plugininstaller.BaseOperation{
			dockerInstaller.HandleRunFailure,
		},
		GetStatusOperation: dockerInstaller.GetStaticStateFromOptions,
	}

	// 3. execute installer get status and error
	rawOptions, err := buildDockerOptions(opts).Encode()
	if err != nil {
		return nil, err
	}
	status, err := runner.Execute(rawOptions)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)

	return status, nil
}
