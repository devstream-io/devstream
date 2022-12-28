package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolTekton = "tekton"

var DefaultConfigWithTekton = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "tekton/tekton-pipeline",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "tekton",
		Namespace:   "tekton",
	},
	Repo: helmCommon.Repo{
		URL:  "https://steinliber.github.io/tekton-helm-chart/",
		Name: "tekton",
	},
}

func init() {
	DefaultOptionsMap[toolTekton] = &DefaultConfigWithTekton
	StatusGetterFuncMap[toolTekton] = GetTektonStatus
}

func GetTektonStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
