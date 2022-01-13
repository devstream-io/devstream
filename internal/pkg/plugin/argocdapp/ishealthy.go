package argocdapp

import argoclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"

func IsHealthy(options *map[string]interface{}) (bool, error) {
	_, _ = argoclient.NewForConfig(nil)
	return true, nil
}
