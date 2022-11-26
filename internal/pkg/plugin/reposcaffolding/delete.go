package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/reposcaffolding"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			reposcaffolding.Validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			reposcaffolding.DeleteRepo,
		},
	}
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}

	// 2. return ture if all process success
	return true, nil
}
