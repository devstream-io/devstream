package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolSonarQube = "sonarqube"

var DefaultConfigWithSonarQube = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "sonarqube/sonarqube",
		Timeout:     "10m",
		Version:     "",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "sonarqube",
		Namespace:   "sonarqube",
	},
	Repo: helmCommon.Repo{
		URL:  "https://SonarSource.github.io/helm-chart-sonarqube",
		Name: "sonarqube",
	},
}

func init() {
	DefaultOptionsMap[toolSonarQube] = &DefaultConfigWithSonarQube
	StatusGetterFuncMap[toolSonarQube] = GetSonarQubeStatus
}

func GetSonarQubeStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
