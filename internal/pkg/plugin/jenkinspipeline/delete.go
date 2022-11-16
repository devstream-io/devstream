package jenkinspipeline

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			setDefault,
			validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			// TODO(daniel-hutao): delete secret: docker-config
			cifile.DeleteCIFiles,
			deletePipeline,
		},
	}
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
