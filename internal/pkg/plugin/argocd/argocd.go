package argocd

import (
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/helm"
)

const (
	defaultRepoName = "argo"
	defaultRepoURL  = "https://argoproj.github.io/argo-helm"
)

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

func defaultMissedOptions(opts *Options) {
	if opts.Repo.URL == "" {
		opts.Repo.URL = defaultRepoURL
	}
	if opts.Repo.Name == "" {
		opts.Repo.Name = defaultRepoName
	}
}
