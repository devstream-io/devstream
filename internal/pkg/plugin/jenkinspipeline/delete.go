package jenkinspipeline

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			jenkins.SetJobDefaultConfig,
			jenkins.ValidateJobConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			// TODO(daniel-hutao): delete secret: docker-config
			ci.DeleteCIFiles,
			jenkins.DeletePipeline,
		},
	}
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
