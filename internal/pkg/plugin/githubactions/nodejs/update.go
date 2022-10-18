package nodejs

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update remove and set up GitHub Actions workflows.
func Update(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			github.Validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			github.ProcessAction(github.ActionUpdate),
		},
		GetStatusOperation: github.GetActionStatus,
	}

	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
