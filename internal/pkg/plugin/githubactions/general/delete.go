package general

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			ci.SetDefault(ciType),
			validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			//TODO(jiafeng meng): delete github secret
			cifile.DeleteCIFiles,
		},
		GetStatusOperation: cifile.GetCIFileStatus,
	}

	// Execute all Operations in Operator
	_, err = operator.Execute(configmanager.RawOptions(options))
	if err != nil {
		return false, err
	}

	return true, nil
}
