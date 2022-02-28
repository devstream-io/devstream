package kubeprometheus

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	. "github.com/merico-dev/stream/internal/pkg/plugin/common/helm"
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create creates kube-prometheus with provided options.
func Create(options *map[string]interface{}) (map[string]interface{}, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return nil, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	if err := DealWithNsWhenInstall(&param); err != nil {
		return nil, err
	}

	var retErr error
	defer func() {
		if retErr == nil {
			return
		}
		if err := DealWithNsWhenInterruption(&param); err != nil {
			log.Errorf("Failed to deal with namespace: %s.", err)
		}
		log.Debugf("Deal with namespace when interruption succeeded.")
	}()

	var h *helm.Helm
	h, retErr = helm.NewHelm(param.GetHelmParam())
	if retErr != nil {
		return nil, retErr
	}

	log.Info("Installing or updating kube-prometheus-stack helm chart ...")
	if retErr = h.InstallOrUpgradeChart(); retErr != nil {
		log.Debugf("Failed to install or upgrade the Chart: %s.", retErr)
		return nil, retErr
	}

	releaseName := param.Chart.ReleaseName
	retMap := GetStaticState(releaseName).ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}
