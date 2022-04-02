package kubeprometheus

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
	KubePrometheusDefaultNamespace = "monitoring"
)

// Read reads the state for kube-prometheus with provided options.
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
		namespace = KubePrometheusDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	retState := &helm.InstanceState{}
	releaseName := opts.Chart.ReleaseName

	err = readDeployments(kubeClient, namespace, releaseName, retState)
	if err != nil {
		log.Debugf("Failed to read deployments: %s.", err)
		return nil, err
	}
	err = readDaemonsets(kubeClient, namespace, releaseName, retState)
	if err != nil {
		log.Debugf("Failed to read daemonsets: %s.", err)
		return nil, err
	}
	err = readStatefulsets(kubeClient, namespace, releaseName, retState)
	if err != nil {
		log.Debugf("Failed to read statefulsets: %s.", err)
		return nil, err
	}

	log.Debugf("All resources read ready.")
	return retState.ToStringInterfaceMap(), nil
}

func readDeployments(kubeClient *k8s.Client, namespace, releaseName string, state *helm.InstanceState) error {
	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		log.Debugf("Failed to list deployments: %s.", err)
		return err
	}

	for _, dp := range dps {
		DefaultDeploymentList := GetDefaultDeploymentList(releaseName)
		dpName := dp.GetName()
		if !slices.Contains(DefaultDeploymentList, dpName) {
			log.Infof("Found unknown deployment: %s.", dpName)
		}

		ready := kubeClient.IsDeploymentReady(&dp)
		state.Workflows.AddDeployment(dpName, ready)
		log.Debugf("The deployment %s is %t.", dp.GetName(), ready)
	}

	return nil
}

func readDaemonsets(kubeClient *k8s.Client, namespace, releaseName string, state *helm.InstanceState) error {
	dss, err := kubeClient.ListDaemonsets(namespace)
	if err != nil {
		log.Debugf("Failed to list daemonsets: %s.", err)
		return err
	}

	for _, ds := range dss {
		DefaultDaemonsetList := GetDefaultDaemonsetList(releaseName)
		dsName := ds.GetName()
		if !slices.Contains(DefaultDaemonsetList, dsName) {
			log.Infof("Found unknown daemonset: %s.", dsName)
		}

		ready := kubeClient.IsDaemonsetReady(&ds)
		state.Workflows.AddDaemonset(dsName, ready)
		log.Debugf("The daemonset %s is %t.", ds.GetName(), ready)
	}

	return nil
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
