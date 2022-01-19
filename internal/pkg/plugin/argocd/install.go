package argocd

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
)

// Install installs ArgoCD with provided options.
func Install(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Printf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	h, err := helm.NewHelm(&param)
	if err != nil {
		return false, err
	}

	log.Println("Installing or updating argocd helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		return false, err
	}

	return true, nil
}
