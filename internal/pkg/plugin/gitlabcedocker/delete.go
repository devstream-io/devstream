package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// 1. create config and pre-handle operations
	opts, err := validateAndDefault(options)
	if err != nil {
		return false, err
	}

	// 2. config delete operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			dockerInstaller.Validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			dockerInstaller.DeleteAll,
		},
	}

	// 3. delete and get status
	rawOptions, err := types.EncodeStruct(buildDockerOptions(opts))
	if err != nil {
		return false, err
	}
	_, err = operator.Execute(rawOptions)
	if err != nil {
		return false, err
	}

	return true, nil
}
