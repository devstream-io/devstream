package localstack

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	CreateNamespace: *types.Bool(false),
	Chart: helmCommon.Chart{
		ChartName:   "localstack/localstack",
		ReleaseName: "localstack",
		Namespace:   "default",
		Timeout:     "5m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
	},
	Repo: helmCommon.Repo{
		URL:  "https://localstack.github.io/helm-charts",
		Name: "localstack-charts",
	},
}
