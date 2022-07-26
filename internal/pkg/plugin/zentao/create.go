package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/goclient"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// 1. config install operations
	runner := &plugininstaller.Runner{
		PreExecuteOperations: []plugininstaller.MutableOperation{
			goclient.Validate,
		},
		ExecuteOperations: []plugininstaller.BaseOperation{
			goclient.DealWithNsWhenInstall,
			goclient.CreatePersistentVolumeWrapper(defaultPVPath),
			goclient.CreatePersistentVolumeClaim,
			goclient.CreateDeploymentWrapperLabelAndContainerPorts(defaultZentaolabels, &defaultZentaoPorts, defaultName),
			goclient.CreateServiceWrapperLabelAndPorts(defaultZentaolabels, &defaultSVCPort),
			goclient.WaitForReady(retryTimes),
		},
		TermateOperations: []plugininstaller.BaseOperation{
			goclient.DealWithErrWhenInstall,
		},
		GetStatusOperation: goclient.GetState,
	}

	// 2. execute installer get status and error
	status, err := runner.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}

	// Currently just store zentao's application status: "running" or "stopped"
	return status, nil
}
