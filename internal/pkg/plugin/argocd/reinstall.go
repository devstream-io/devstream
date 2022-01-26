package argocd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
)

func Reinstall(options *map[string]interface{}) (bool, error) {
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

	log.Info("Uninstalling argocd helm chart ...")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	// delete the namespace
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}
	if err = kubeClient.DeleteNamespace(param.Chart.Namespace); err != nil {
		log.Errorf("Failed to delete the %s namespace: %s", param.Chart.Namespace, err)
		return false, err
	}

	// install
	log.Info("Installing or updating argocd helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		return false, err
	}

	return true, nil
}
