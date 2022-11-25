package devlakeconfig

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			validate,
			RenderAuthConfig,
		},
		ExecuteOperations: installer.ExecuteOperations{
			DeleteConfig,
		},
	}

	// Execute all Operations in Operator
	_, err = operator.Execute(options)
	if err != nil {
		return false, err
	}

	return true, nil
}
