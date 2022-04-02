package argocd

import "github.com/devstream-io/devstream/pkg/util/helm"

var DefaultDeploymentList = []string{
	"argocd-applicationset-controller",
	"argocd-dex-server",
	"argocd-notifications-controller",
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
