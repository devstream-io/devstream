package python

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update remove and set up GitHub Actions workflows.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			github.Validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			deleteDockerHubInfoForPush,
			createDockerHubInfoForPush,
			github.ProcessAction("update"),
		},
		GetStatusOperation: github.GetActionState,
	}

	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
