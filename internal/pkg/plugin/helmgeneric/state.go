package helmgeneric

import "github.com/devstream-io/devstream/pkg/util/helm"

func GetStaticState(DefaultDeploymentList []string) *helm.InstanceState {
	retState := &helm.InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}
