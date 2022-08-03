package argocd

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "argo/argo-cd",
		Timeout:     "5m",
		UpgradeCRDs: helmCommon.GetBoolTrueAddress(),
		Wait:        helmCommon.GetBoolTrueAddress(),
	},
	Repo: helmCommon.Repo{
		URL:  "https://argoproj.github.io/argo-helm",
		Name: "argo",
	},
}
