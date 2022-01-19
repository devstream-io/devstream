package kubeprometheus

import "github.com/merico-dev/stream/internal/pkg/util/helm"

func validate(param *helm.HelmParam) []error {
	return helm.Validate(param)
}
