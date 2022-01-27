package kubeprometheus

import "github.com/merico-dev/stream/pkg/util/helm"

// Param is the struct for parameters used by the kubeprometheus package.
type Param struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	helm.HelmParam
}
