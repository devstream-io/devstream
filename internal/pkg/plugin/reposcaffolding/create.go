package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/reposcaffolding"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			reposcaffolding.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			reposcaffolding.InstallRepo,
		},
		GetStatusOperation: reposcaffolding.GetDynamicStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
