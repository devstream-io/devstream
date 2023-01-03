package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

const toolHarbor = "harbor"

var DefaultConfigWithHarbor = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "harbor/harbor",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "harbor",
		Namespace:   "harbor",
	},
	Repo: helmCommon.Repo{
		URL:  "https://helm.goharbor.io",
		Name: "harbor",
	},
}

func init() {
	RegisterDefaultHelmAppInstance(toolHarbor, &DefaultConfigWithHarbor, GetHarborStatus)
}

func GetHarborStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
