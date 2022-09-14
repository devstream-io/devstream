package artifactory

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "jfrog/artifactory",
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
