package helm

import (
	"github.com/devstream-io/devstream/pkg/util/helm"
)

// NOTICE: Don't use:
// type Param struct {
// 	CreateNamespace bool `mapstructure:"create_namespace"`
// 	helm.HelmParam
// }
// or
// type Param struct {
// 	CreateNamespace bool `mapstructure:"create_namespace"`
// 	*helm.HelmParam
// }
// see pr #174 for more info

// Param is the struct for parameters used by the argocd package.
type Options struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	Repo            helm.Repo
	Chart           helm.Chart
}
