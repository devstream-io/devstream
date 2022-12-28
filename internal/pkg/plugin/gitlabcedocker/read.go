package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// 1. create config and pre-handle operations
	opts, err := validateAndDefault(options)
	if err != nil {
		return nil, err
	}

	// 2. config read operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			dockerInstaller.Validate,
		},
		GetStatusOperation: dockerInstaller.GetRunningStatus,
	}

	// 3. get status
	rawOptions, err := types.EncodeStruct(buildDockerOptions(opts))
	if err != nil {
		return nil, err
	}
	status, err := operator.Execute(rawOptions)
	if err != nil {
		return nil, err
	}

	log.Debugf("Return map: %v", status)
	return status, nil
}
