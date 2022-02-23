package kubeprometheus

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
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

	if err := dealWithNsWhenInstall(&param); err != nil {
		return nil, err
	}

	h, err := helm.NewHelm(param.GetHelmParam())
	if err != nil {
		return nil, err
	}

	log.Info("Installing or updating kube-prometheus-stack helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		log.Debugf("Failed to install or upgrade the Chart: %s.", err)
		if err = dealWithNsWhenInterruption(&param); err != nil {
			log.Debugf("Failed to deal with namespace: %s.", err)
			// don't need to return this err here, just print it.
			// The err return by InstallOrUpgradeChart() is more useful for the caller.
		}
		return nil, err
	}

	releaseName := param.Chart.ReleaseName
	retMap := GetStaticState(releaseName).ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}

// TODO(daniel-hutao): All helm-style plugins has this code logic, maybe it's better to move it to a common package.
func dealWithNsWhenInstall(param *Param) error {
	if !param.CreateNamespace {
		return nil
	}

	log.Debugf("Prepare to create the namespace: %s.", param.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.CreateNamespace(param.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to create the namespace: %s.", param.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been created.", param.Chart.Namespace)
	return nil
}

// TODO(daniel-hutao): All helm-style plugins has this code logic, maybe it's better to move it to a common package.
func dealWithNsWhenInterruption(param *Param) error {
	if !param.CreateNamespace {
		return nil
	}

	log.Debugf("Prepare to delete the namespace: %s.", param.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.DeleteNamespace(param.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to delete the namespace: %s.", param.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been deleted.", param.Chart.Namespace)
	return nil
}
