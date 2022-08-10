package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "argo/argo-cd",
		Timeout:     "5m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "argocd",
	},
	CreateNamespace: types.Bool(false),
	Repo: helmCommon.Repo{
		URL:  "https://argoproj.github.io/argo-helm",
		Name: "argo",
	},
}
