package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		GetStatusOperation: github.GetActionState,
	}

	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}
	return status, nil
}
