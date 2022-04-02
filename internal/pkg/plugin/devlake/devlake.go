package devlake

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const devLakeInstallYAMLDownloadURL = "https://raw.githubusercontent.com/merico-dev/lake/main/k8s-deploy.yaml"
const devLakeInstallYAMLFileName = "devlake-k8s-deploy.yaml"

// according to devLakeInstallYAMLDownloadURL
// a successful DevLake installation should have the following deployments
// (and corresponding services as well)
var devLakeDeployments = [4]string{
	"mysql",
	"grafana",
	"config-ui",
	"devlake",
}

func buildState(opts Options) map[string]interface{} {
	res := make(map[string]interface{})

	res["deployments"] = make(map[string]interface{})
	res["services"] = make(map[string]interface{})
	for _, d := range devLakeDeployments {
		res["deployments"].(map[string]interface{})[d] = true
		res["services"].(map[string]interface{})[d] = true
	}

	return res
}

func allDeploymentsAndServicesReady() error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// TODO(ironcore864): now the namespace is hard-coded instead of parsed from the YAML file
	namespace := "devlake"

	// check if all deployments are ready
	for _, d := range devLakeDeployments {
		dp, err := kubeClient.GetDeployment(namespace, d)
		if err != nil {
			return err
		}

		if kubeClient.IsDeploymentReady(dp) {
			log.Infof("The deployment %s is ready.", dp.Name)
			continue
		} else {
			log.Infof("The deployment %s is not ready yet.", dp.Name)
			return fmt.Errorf("deployment %s not ready", dp.Name)
		}
	}

	// check if all services exist
	for _, d := range devLakeDeployments {
		_, err := kubeClient.GetService(namespace, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func readDeploymentsAndServicesAndBuildState() (map[string]interface{}, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	// TODO(ironcore864): now the namespace is hard-coded instead of parsed from the YAML file
	namespace := "devlake"

	res := make(map[string]interface{})

	res["deployments"] = make(map[string]interface{})
	res["services"] = make(map[string]interface{})

	// check if all deployments are ready
	for _, d := range devLakeDeployments {
		// deployment
		dp, err := kubeClient.GetDeployment(namespace, d)
		if err == nil && kubeClient.IsDeploymentReady(dp) {
			res["deployments"].(map[string]interface{})[d] = true
		} else {
			res["deployments"].(map[string]interface{})[d] = true
		}

		// services
		_, err = kubeClient.GetService(namespace, d)
		if err == nil {
			res["services"].(map[string]interface{})[d] = true
		} else {
			res["services"].(map[string]interface{})[d] = false
		}
	}

	log.Debugf("Resource read returns: %v.", res)
	return res, nil
}
