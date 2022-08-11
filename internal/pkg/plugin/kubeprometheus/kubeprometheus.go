package kubeprometheus

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "prometheus-community/kube-prometheus-stack",
		Timeout:     "5m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "prometheus",
		Namespace:   "prometheus",
	},
	CreateNamespace: types.Bool(false),
	Repo: helmCommon.Repo{
		URL:  "https://prometheus-community.github.io/helm-charts",
		Name: "prometheus-community",
	},
}
