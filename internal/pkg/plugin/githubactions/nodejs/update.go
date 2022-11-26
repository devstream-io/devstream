package nodejs

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/github"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update remove and set up GitHub Actions workflows.
func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			github.Validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: installer.ExecuteOperations{
			github.ProcessAction(github.ActionUpdate),
		},
		GetStatusOperation: github.GetActionStatus,
	}

	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
