package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/github"
)

// Delete remove GitHub Actions workflows.
func Delete(options configmanager.RawOptions) (bool, error) {
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: installer.ExecuteOperations{
			deleteDockerHubInfoForPush,
			github.ProcessAction(github.ActionDelete),
		},
	}

	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}
