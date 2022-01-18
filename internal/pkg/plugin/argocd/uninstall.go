package argocd

import (
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
)

func Uninstall(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	h, err := helm.NewHelm(&param)
	if err != nil {
		return false, err
	}

	log.Println("uninstalling argocd helm chart")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	return true, nil
}
