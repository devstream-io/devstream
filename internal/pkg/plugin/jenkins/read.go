package jenkins

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
	JenkinsDefaultNamespace = "jenkins"
)

// Read reads the state for jenkins with provided options.
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
		namespace = JenkinsDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	retState := &helm.InstanceState{}
	releaseName := opts.Chart.ReleaseName

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
