package argocd

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/util/k8s"
)

// The deployments should exist:
// 1. argocd-application-controller
// 2. argocd-dex-server
// 3. argocd-redis
// 4. argocd-repo-server
// 5. argocd-server
const (
	DefaultDeploymentCount = 5
	ArgoCDNamespace        = "argocd"
)

func IsHealthy(options *map[string]interface{}) (bool, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}

	dps, err := kubeClient.ListDeployments(ArgoCDNamespace)
	if err != nil {
		return false, err
	}

	if len(dps) != DefaultDeploymentCount {
		return false, fmt.Errorf("expect deployments count %d, but got count %d now", DefaultDeploymentCount, len(dps))
	}

	hasNotReadyDeployment := false
	for _, dp := range dps {
		if kubeClient.IsDeploymentReady(&dp) {
			log.Printf("the deployment %s is ready", dp.Name)
			continue
		}
		log.Printf("the deployment %s is not ready", dp.Name)
		hasNotReadyDeployment = true
	}

	if hasNotReadyDeployment {
		return false, fmt.Errorf("some deployment are not ready")
	}
	return true, nil
}
