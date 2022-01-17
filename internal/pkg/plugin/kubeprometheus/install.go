package kubeprometheus

import (
	"github.com/merico-dev/stream/internal/pkg/util/helm"
	"github.com/mitchellh/mapstructure"
)

// Install installs kube-prometheus with provided options.
func Install(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	h, err := helm.NewHelm(&param)
	if err != nil {
		return false, err
	}

	if err = h.InstallOrUpgradeChart(); err != nil{
		return false, err
	}

	return true, nil
}
