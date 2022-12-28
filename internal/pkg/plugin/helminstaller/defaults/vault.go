package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolVault = "vault"

var DefaultConfigWithVault = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "hashicorp/vault",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "vault",
		Namespace:   "vault",
	},
	Repo: helmCommon.Repo{
		URL:  "https://helm.releases.hashicorp.com",
		Name: "hashicorp",
	},
}

func init() {
	DefaultOptionsMap[toolVault] = &DefaultConfigWithVault
	StatusGetterFuncMap[toolVault] = GetVaultStatus
}

func GetVaultStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
