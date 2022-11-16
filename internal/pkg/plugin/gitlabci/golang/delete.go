package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			cifile.SetDefaultConfig(gitlabci.DefaultCIOptions),
			setCIContent,
			cifile.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			cifile.DeleteCIFiles,
		},
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}
