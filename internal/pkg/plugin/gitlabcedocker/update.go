package gitlabcedocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// 1. create config and pre-handle operations
	opts, err := validateAndDefault(options)
	if err != nil {
		return nil, err
	}

	// reserve data when updated
	opts.RmDataAfterDelete = types.Bool(false)

	// 2. config install operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			dockerInstaller.Validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			dockerInstaller.DeleteAll,
			dockerInstaller.Install,
			showHelpMsg,
		},
		GetStatusOperation: dockerInstaller.GetRunningStatus,
	}

	// 3. update and get status
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
