package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

// TODO(daniel-hutao): update the config below after devlake chart released.
var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "devlake/devlake",
		Timeout:     "5m",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "devlake",
		Namespace:   "devlake",
	},
	Repo: helmCommon.Repo{
		URL:  "https://merico-dev.github.io/devlake-helm-chart",
		Name: "devlake",
	},
}
