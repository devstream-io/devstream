package kubeprometheus

import "github.com/merico-dev/stream/pkg/util/helm"

// validate validates the options provided by the core.
func validate(param *Param) []error {
	return helm.Validate(param.GetHelmParam())
}
