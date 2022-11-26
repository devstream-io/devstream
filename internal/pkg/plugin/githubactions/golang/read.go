package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/github"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		GetStatusOperation: github.GetActionStatus,
	}

	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	return status, nil
}
