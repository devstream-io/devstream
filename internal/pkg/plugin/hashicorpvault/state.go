package hashicorpvault

import "github.com/devstream-io/devstream/pkg/util/helm"

var DefaultDeploymentList = []string{
	"vault-agent-injector",
}

func GetStaticState() *helm.InstanceState {
	retState := &helm.InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}
