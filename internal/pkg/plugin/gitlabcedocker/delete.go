package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. create config and pre-handle operations
	opts, err := preHandleOptions(options)
	if err != nil {
		return false, err
	}

	// 2. config delete operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			dockerInstaller.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			dockerInstaller.Delete,
		},
	}

	// 3. delete and get status
	rawOptions, err := buildDockerOptions(opts).Encode()
	if err != nil {
		return false, err
	}
	_, err = runner.Execute(rawOptions)
	if err != nil {
		return false, err
	}

	return true, nil
}
