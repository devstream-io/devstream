package argocd

import "github.com/merico-dev/stream/pkg/util/helm"

// Param is the struct for parameters used by the argocd package.
type Param struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	helm.HelmParam
}
