package sonarqube

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "sonarqube/sonarqube",
		Timeout:     "20m",
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
