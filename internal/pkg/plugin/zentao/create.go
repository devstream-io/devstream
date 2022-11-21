package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/goclient"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &installer.Operator{
		PreExecuteOperations: installer.PreExecuteOperations{
			goclient.Validate,
		},
		ExecuteOperations: installer.ExecuteOperations{
			goclient.DealWithNsWhenInstall,
			goclient.CreatePersistentVolumeWrapper(defaultPVPath),
			goclient.CreatePersistentVolumeClaim,
			goclient.CreateDeploymentWrapperLabelAndContainerPorts(defaultZentaolabels, &defaultZentaoPorts, defaultName),
			goclient.CreateServiceWrapperLabelAndPorts(defaultZentaolabels, &defaultSVCPort),
			goclient.WaitForReady(retryTimes),
		},
		TerminateOperations: installer.TerminateOperations{
			goclient.DealWithErrWhenInstall,
		},
		GetStatusOperation: goclient.GetStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}

	// Currently just store zentao's application status: "running" or "stopped"
	return status, nil
}
