package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			renderConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			dockerInstaller.ComposeDown,
			dockerInstaller.ComposeUp,
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
