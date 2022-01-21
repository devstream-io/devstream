package argocd

import (
	"fmt"
	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
)

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

	log.Info("uninstalling argocd helm chart")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	return true, nil
}
