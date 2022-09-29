package apachedevlake

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

// TODO(daniel-hutao): update the config below after devlake chart released.
var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "apache-devlake/devlake",
		Timeout:     "5m",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "devlake",
		Namespace:   "devlake",
	},
	Repo: helmCommon.Repo{
		URL:  "https://apache-devlake.github.io/apache-devlake-helm",
		Name: "apache-devlake",
	},
}
