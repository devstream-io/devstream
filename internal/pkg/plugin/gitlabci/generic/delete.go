package generic

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
)

func Delete(options map[string]interface{}) (bool, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			ci.SetDefaultConfig(gitlabci.DefaultCIOptions),
			ci.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			ci.DeleteCIFiles,
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}
	return true, nil
}
