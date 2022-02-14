package devlake

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/k8s"
)

const (
	DevLakeTotalK8sDeployments = 4
)

func IsHealthy(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	if errs := validateParams(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}

	// TODO(ironcore864): now the namespace is hard-coded instead of parsed from the YAML file
	namespace := "devlake"
	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		return false, err
	}

	// check if the number of deployments is correct
	if len(dps) != DevLakeTotalK8sDeployments {
		return false, fmt.Errorf("expect %d deployments, but only got %d",
			DevLakeTotalK8sDeployments, len(dps))
	}

	// check if all deployments are ready
	deploymentNotReady := false
	for _, dp := range dps {
		if kubeClient.IsDeploymentReady(&dp) {
			log.Infof("Deployment %s is ready.", dp.Name)
			continue
		}
		deploymentNotReady = true
		log.Infof("Deployment %s is not ready.", dp.Name)
		break
	}
	if deploymentNotReady {
		return false, fmt.Errorf("some deployment(s) is not ready")
	}

	return true, nil
}
