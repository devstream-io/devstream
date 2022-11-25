package helminstaller

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			RenderDefaultConfig,
			RenderValuesYaml,
			validate,
		},
		ExecuteOperations:   helm.DefaultUpdateOperations,
		TerminateOperations: helm.DefaultTerminateOperations,
		GetStatusOperation:  IndexStatusGetterFunc(options),
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v", status)
	return status, nil
}
