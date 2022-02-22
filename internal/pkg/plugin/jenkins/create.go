package jenkins

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
)

// Create creates jenkins with provided options.
func Create(options *map[string]interface{}) (map[string]interface{}, error) {
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

	param.renderValuesYamlForJenkins()

	if err := dealWithNsWhenInstall(&param); err != nil {
		return nil, err
	}

	h, err := helm.NewHelm(param.GetHelmParam())
	if err != nil {
		return nil, err
	}

	// pre-create
	if err = preCreate(); err != nil {
		log.Errorf("The pre-create logic failed: %s", err)
		return nil, err
	}

	log.Info("Installing or updating jenkins helm chart ...")
	if err = h.InstallOrUpgradeChart(); err != nil {
		return nil, err
	}

	releaseName := param.Chart.ReleaseName
	retMap := GetStaticState(releaseName).ToStringInterfaceMap()
	log.Debugf("Return map: %v", retMap)

	return retMap, nil
}

func dealWithNsWhenInstall(param *Param) error {
	if !param.CreateNamespace {
		return nil
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	return kubeClient.CreateNamespace(param.Chart.Namespace)
}
