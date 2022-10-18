package harbordocker

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (statemanager.ResourceStatus, error) {
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
	state, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", state)
	return state, nil
}
