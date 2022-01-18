package argocd

import (
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
)

func Reinstall(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	h, err := helm.NewHelm(&param)
	if err != nil {
		return false, err
	}

	log.Println("Uninstalling argocd helm chart ...")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	log.Println("Installing or updating argocd helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		return false, err
	}

	return true, nil
}
