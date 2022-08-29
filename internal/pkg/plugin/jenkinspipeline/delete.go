package jenkinspipeline

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			jenkins.SetJobDefaultConfig,
			jenkins.ValidateJobConfig,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			ci.DeleteCIFiles,
			jenkins.DeleteJob,
		},
	}
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
