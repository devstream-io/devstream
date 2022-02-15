package argocd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/helm"
)

func Update(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	h, err := helm.NewHelm(param.GetHelmParam())
	if err != nil {
		return nil, err
	}

	log.Info("Uninstalling argocd helm chart ...")
	if err = h.UninstallHelmChartRelease(); err != nil {
		return nil, err
	}

	if err := dealWithNsWhenUninstall(&param); err != nil {
		return nil, err
	}

	// install
	if err := dealWithNsWhenInstall(&param); err != nil {
		return nil, err
	}

	log.Info("Installing or updating argocd helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		return nil, err
	}

	return make(map[string]interface{}), nil
}
