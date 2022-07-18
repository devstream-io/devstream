package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

var (
	defaultDeploymentList = []string{
		"argocd-applicationset-controller",
		"argocd-dex-server",
		"argocd-notifications-controller",
		"argocd-redis",
		"argocd-repo-server",
		"argocd-server"}
	defaultRepoURL  = "https://argoproj.github.io/argo-helm"
	defaultRepoName = "argo"
)

func defaultMissedOption(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	if opts.Repo.URL == "" {
		opts.Repo.URL = defaultRepoURL
	}
	if opts.Repo.Name == "" {
		opts.Repo.Name = defaultRepoName
	}
	return opts.Encode()
}
