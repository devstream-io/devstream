package kubeprometheus

import "github.com/merico-dev/stream/pkg/util/helm"

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

// Param is the struct for parameters used by the kubeprometheus package.
type Param struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	Repo            helm.Repo
	Chart           helm.Chart
}

func (p *Param) GetHelmParam() *helm.HelmParam {
	return &helm.HelmParam{
		Repo:  p.Repo,
		Chart: p.Chart,
	}
}
