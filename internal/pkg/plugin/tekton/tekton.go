package tekton

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "tekton/tekton-pipeline",
		Timeout:     "5m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "tekton",
		Namespace:   "tekton",
	},
	CreateNamespace: types.Bool(true),
	Repo: helmCommon.Repo{
		URL:  "https://steinliber.github.io/tekton-helm-chart/",
		Name: "tekton",
	},
}
