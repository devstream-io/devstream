package argocd

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"k8s.io/utils/strings/slices"

	. "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	ArgocdDefaultNamespace = "argocd"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	namespace := opts.Chart.Namespace
	if namespace == "" {
		namespace = ArgocdDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		return nil, err
	}

	retState := &helm.InstanceState{}
	for _, dp := range dps {
		dpName := dp.GetName()
		if !slices.Contains(DefaultDeploymentList, dpName) {
			log.Infof("Found unknown deployment: %s.", dpName)
		}

		ready := kubeClient.IsDeploymentReady(&dp)
		retState.Workflows.AddDeployment(dpName, ready)
		log.Debugf("The deployment %s is %t.", dp.GetName(), ready)
	}

	retMap := retState.ToStringInterfaceMap()
	log.Debugf("Return map: %v.", retMap)

	return retMap, nil
}
