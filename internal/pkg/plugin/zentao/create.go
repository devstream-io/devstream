package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/goclient"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			goclient.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			goclient.DealWithNsWhenInstall,
			goclient.CreatePersistentVolumeWrapper(defaultPVPath),
			goclient.CreatePersistentVolumeClaim,
			goclient.CreateDeploymentWrapperLabelAndContainerPorts(defaultZentaolabels, &defaultZentaoPorts, defaultName),
			goclient.CreateServiceWrapperLabelAndPorts(defaultZentaolabels, &defaultSVCPort),
			goclient.WaitForReady(retryTimes),
		},
		TerminateOperations: plugininstaller.TerminateOperations{
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
