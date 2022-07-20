package util

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func CheckAllDeployAndServiceReady(namespace string, deployList []string) error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// check if all deployments are ready
	for _, d := range deployList {
		dp, err := kubeClient.GetDeployment(namespace, d)
		if err != nil {
			return err
		}

		if !kubeClient.IsDeploymentReady(dp) {
			log.Infof("The deployment %s is not ready yet.", dp.Name)
			return fmt.Errorf("deployment %s not ready", dp.Name)
		}
		log.Infof("The deployment %s is ready.", dp.Name)
	}

	// check if all services exist
	for _, d := range deployList {
		_, err := kubeClient.GetService(namespace, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadDepAndServiceState(namespace string, deployList []string) (map[string]interface{}, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	res["deployments"] = make(map[string]interface{})
	res["services"] = make(map[string]interface{})

	// check if all deployments are ready
	for _, d := range deployList {
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

func GetArgoCDAppFromK8sAndSetState(state map[string]interface{}, name, namespace string) error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	app, err := kubeClient.GetArgocdApplication(namespace, name)
	if err != nil {
		return err
	}

	d := kubeClient.DescribeArgocdApp(app)
	state["app"] = d["app"]
	state["src"] = d["src"]
	state["dest"] = d["dest"]

	return nil
}
