package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			renderConfig,
		},
		ExecuteOperations: installer.ExecuteOperations{
			Install,
		},
		GetStatusOperation: dockerInstaller.ComposeStatus,
	}

	// Execute all Operations in Operator
	state, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", state)
	return state, nil
}
