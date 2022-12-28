package defaults

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolOpenLDAP = "openldap"

var DefaultConfigWithOpenLDAP = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "helm-openldap/openldap-stack-ha",
		Version:     "",
		Timeout:     "10m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "openldap",
		Namespace:   "openldap",
	},
	Repo: helmCommon.Repo{
		URL:  "https://jp-gouin.github.io/helm-openldap/",
		Name: "helm-openldap",
	},
}

func init() {
	DefaultOptionsMap[toolOpenLDAP] = &DefaultConfigWithOpenLDAP
	StatusGetterFuncMap[toolOpenLDAP] = GetOpenLDAPStatus
}

func GetOpenLDAPStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return helm.GetAllResourcesStatus(options)
}
