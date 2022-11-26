package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			renderConfig,
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
