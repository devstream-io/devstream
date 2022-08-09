package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
)

func Delete(options map[string]interface{}) (bool, error) {
	// 1. create config and pre-handle operations
	opts, err := validateAndDefault(options)
	if err != nil {
		return false, err
	}

	// 2. config delete operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			dockerInstaller.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			dockerInstaller.DeleteAll,
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
