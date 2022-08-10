package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create sets up GitHub Actions workflow(s).
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			createDockerHubInfoForPush,
			github.ProcessAction("create"),
		},
		GetStateOperation: github.GetStaticWorkFlowState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
