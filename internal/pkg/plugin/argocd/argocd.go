package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "argo/argo-cd",
		Timeout:     "5m",
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
