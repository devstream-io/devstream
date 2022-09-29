package apachedevlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func Delete(options map[string]interface{}) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			// TODO(dtm): Add your PreExecuteOperations here.
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			// TODO(dtm): Add your ExecuteOperations here.
		},
	}
	_, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
