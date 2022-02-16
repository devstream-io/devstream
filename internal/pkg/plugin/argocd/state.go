package argocd

import "github.com/merico-dev/stream/pkg/util/helm"

var DefaultDeploymentList = []string{
	"argocd-application-controller",
	"argocd-dex-server",
	"argocd-redis",
	"argocd-repo-server",
	"argocd-server",
}

func GetStaticState() *helm.InstanceState {
	retState := &helm.InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}
