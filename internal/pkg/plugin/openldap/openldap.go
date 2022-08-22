package openldap

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "helm-openldap/openldap-stack-ha",
		Timeout:     "5m",
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
