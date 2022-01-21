package argocd

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
	"github.com/merico-dev/stream/internal/pkg/util/k8s"
)

// The deployments should exist:
// 1. argocd-application-controller
// 2. argocd-dex-server
// 3. argocd-redis
// 4. argocd-repo-server
// 5. argocd-server
const (
	ArgocdDefaultDeploymentCount = 5
	ArgocdDefaultNamespace       = "argocd"
)

func IsHealthy(options *map[string]interface{}) (bool, error) {
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

	namespace := param.Chart.Namespace
	if namespace == "" {
		namespace = ArgocdDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}

	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		return false, err
	}

	if len(dps) != ArgocdDefaultDeploymentCount {
		return false, fmt.Errorf("expect deployments count %d, but got count %d now",
			ArgocdDefaultDeploymentCount, len(dps))
	}

	hasNotReadyDeployment := false
	for _, dp := range dps {
		if kubeClient.IsDeploymentReady(&dp) {
			log.Infof("the deployment %s is ready", dp.Name)
			continue
		}
		log.Infof("the deployment %s is not ready", dp.Name)
		hasNotReadyDeployment = true
	}

	if hasNotReadyDeployment {
		return false, fmt.Errorf("some deployments are not ready")
	}
	return true, nil
}
