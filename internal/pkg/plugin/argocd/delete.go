package argocd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
	"github.com/merico-dev/stream/pkg/util/log"
)

func Delete(options *map[string]interface{}) (bool, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	h, err := helm.NewHelm(param.GetHelmParam())
	if err != nil {
		return false, err
	}

	log.Info("Uninstalling argocd helm chart.")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return false, err
	}

	if err := dealWithNsWhenUninstall(&param); err != nil {
		return false, err
	}

	return true, nil
}

func dealWithNsWhenUninstall(param *Param) error {
	if !param.CreateNamespace {
		return nil
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	return kubeClient.DeleteNamespace(param.Chart.Namespace)
}
