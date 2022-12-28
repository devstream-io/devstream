package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolArgoCD = "argocd"

var DefaultConfigWithArgoCD = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "argo/argo-cd",
		Version:     "",
		Timeout:     "10m",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "argocd",
		Namespace:   "argocd",
	},
	Repo: helmCommon.Repo{
		URL:  "https://argoproj.github.io/argo-helm",
		Name: "argo",
	},
}

func init() {
	DefaultOptionsMap[toolArgoCD] = &DefaultConfigWithArgoCD
	StatusGetterFuncMap[toolArgoCD] = GetArgoCDStatus
}

func GetArgoCDStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
