package argocdapp

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/k8s"
)

func IsHealthy(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}

	namespace := param.App.Namespace
	name := param.App.Name
	app, err := kubeClient.GetArgocdApplication(namespace, name)
	if err != nil {
		return false, err
	}

	if kubeClient.IsArgocdApplicationReady(app) {
		log.Printf("%s/%s is ready", namespace, name)
		return true, nil
	}
	return false, fmt.Errorf("the Argocd Application %s/%s is not ready", namespace, name)
}
