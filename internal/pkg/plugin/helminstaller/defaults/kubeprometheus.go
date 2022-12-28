package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolKubePrometheus = "kube-prometheus"

var DefaultConfigWithKubePrometheus = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "prometheus-community/kube-prometheus-stack",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "prometheus",
		Namespace:   "prometheus",
	},
	Repo: helmCommon.Repo{
		URL:  "https://prometheus-community.github.io/helm-charts",
		Name: "prometheus-community",
	},
}

func init() {
	DefaultOptionsMap[toolKubePrometheus] = &DefaultConfigWithKubePrometheus
	StatusGetterFuncMap[toolKubePrometheus] = GetKubePrometheusStatus
}

func GetKubePrometheusStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
