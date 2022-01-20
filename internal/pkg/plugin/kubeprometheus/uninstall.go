package kubeprometheus

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
	"github.com/merico-dev/stream/internal/pkg/util/log"
)

// Uninstall uninstalls kube-prometheus with provided options.
func Uninstall(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	h, err := helm.NewHelm(&param)
	if err != nil {
		return false, err
	}

	log.Info("Uninstalling kube-prometheus-stack helm chart ...")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	return true, nil
}
