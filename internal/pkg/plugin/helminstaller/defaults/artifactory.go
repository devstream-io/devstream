package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolArtifactory = "artifactory"

var DefaultConfigWithArtifactory = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "jfrog/artifactory",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "artifactory",
		Namespace:   "artifactory",
	},
	Repo: helmCommon.Repo{
		URL:  "https://charts.jfrog.io",
		Name: "jfrog",
	},
}

func init() {
	DefaultOptionsMap[toolArtifactory] = &DefaultConfigWithArtifactory
	StatusGetterFuncMap[toolArtifactory] = GetArtifactoryStatus
}

func GetArtifactoryStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
