package jenkins

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"k8s.io/utils/strings/slices"

	. "github.com/merico-dev/stream/internal/pkg/plugin/common/helm"
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
	"github.com/merico-dev/stream/pkg/util/log"
)

const (
	JenkinsDefaultNamespace = "jenkins"
)

// Read reads the state for jenkins with provided options.
func Read(options *map[string]interface{}) (map[string]interface{}, error) {
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

	namespace := param.Chart.Namespace
	if namespace == "" {
		namespace = JenkinsDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	retState := &helm.InstanceState{}
	releaseName := param.Chart.ReleaseName

	err = readStatefulsets(kubeClient, namespace, releaseName, retState)
	if err != nil {
		log.Debugf("Failed to read statefulsets: %s.", err)
		return nil, err
	}

	log.Debugf("All resources read ready.")
	return retState.ToStringInterfaceMap(), nil
}

func readStatefulsets(kubeClient *k8s.Client, namespace, releaseName string, state *helm.InstanceState) error {
	sss, err := kubeClient.ListStatefulsets(namespace)
	if err != nil {
		log.Debugf("Failed to list statefulsets: %s.", err)
		return err
	}

	for _, ss := range sss {
		DefaultStatefulsetList := GetDefaultStatefulsetList(releaseName)
		ssName := ss.GetName()
		if !slices.Contains(DefaultStatefulsetList, ssName) {
			log.Infof("Found unknown statefulset: %s.", ssName)
		}

		ready := kubeClient.IsStatefulsetReady(&ss)
		state.Workflows.AddStatefulset(ssName, ready)
		log.Debugf("The statefulset %s is %t.", ss.GetName(), ready)
	}

	return nil
}
