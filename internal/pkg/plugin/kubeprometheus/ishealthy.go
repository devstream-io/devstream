package kubeprometheus

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/helm"
	"github.com/merico-dev/stream/internal/pkg/util/k8s"
)

const (
	KubePrometheusDefaultNamespace = "monitoring"

	// The deployments should exist:
	// 1. ${release-name}-grafana
	// 2. ${release-name}-kube-prometheus-stack-operator
	// 3. ${release-name}-kube-state-metrics
	KubePrometheusDefaultDeploymentCount = 3

	// // The statefulsets should exist:
	// 1. alertmanager-${release-name}-kube-prometheus-stack-alertmanager
	// 2. prometheus-${release-name}-kube-prometheus-stack-prometheus
	KubePrometheusDefaultStatefulsetCount = 2

	// // The daemonsets should exist:
	// 1. ${release-name}-prometheus-node-exporter
	KubePrometheusDefaultDaemonsetCount = 1
)

// IsHealthy check the health for kube-prometheus with provided options.
func IsHealthy(options *map[string]interface{}) (bool, error) {
	var param helm.HelmParam
	if err := mapstructure.Decode(*options, &param); err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Printf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	namespace := param.Chart.Namespace
	if namespace == "" {
		namespace = KubePrometheusDefaultNamespace
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return false, err
	}

	hasSomeResorucesNotReady := false

	dpReady, err := isDeploymentsReady(kubeClient, namespace)
	if err != nil {
		return false, err
	}
	if !dpReady {
		hasSomeResorucesNotReady = true
		log.Printf("Some deployments are not ready.")
	}

	dsReady, err := isDaemonsetsReady(kubeClient, namespace)
	if err != nil {
		return false, err
	}
	if !dsReady {
		hasSomeResorucesNotReady = true
		log.Printf("Some daemonsets are not ready.")
	}

	ssReady, err := isStatefulsetsReady(kubeClient, namespace)
	if err != nil {
		return false, err
	}
	if !ssReady {
		hasSomeResorucesNotReady = true
		log.Printf("Some statefulsets are not ready.")
	}

	if hasSomeResorucesNotReady {
		return false, fmt.Errorf("some resources are not ready")
	}

	return true, nil
}

func isDeploymentsReady(kubeClient *k8s.Client, namespace string) (bool, error) {
	dps, err := kubeClient.ListDeployments(namespace)
	if err != nil {
		return false, err
	}

	if len(dps) != KubePrometheusDefaultDeploymentCount {
		return false, fmt.Errorf("expect deployments count %d, but got count %d now",
			KubePrometheusDefaultDeploymentCount, len(dps))
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
		return false, fmt.Errorf("some deployments are not ready")
	}
	return true, nil
}

func isDaemonsetsReady(kubeClient *k8s.Client, namespace string) (bool, error) {
	dss, err := kubeClient.ListDaemonsets(namespace)
	if err != nil {
		return false, err
	}

	if len(dss) != KubePrometheusDefaultDaemonsetCount {
		return false, fmt.Errorf("expect daemonsets count %d, but got count %d now",
			KubePrometheusDefaultDaemonsetCount, len(dss))
	}

	hasNotReadyDaemonset := false
	for _, ds := range dss {
		if kubeClient.IsDaemonsetReady(&ds) {
			log.Printf("the daemonset %s is ready", ds.Name)
			continue
		}
		log.Printf("the daemonset %s is not ready", ds.Name)
		hasNotReadyDaemonset = true
	}

	if hasNotReadyDaemonset {
		return false, fmt.Errorf("some daemonsets are not ready")
	}
	return true, nil
}

func isStatefulsetsReady(kubeClient *k8s.Client, namespace string) (bool, error) {
	sss, err := kubeClient.ListStatefulsets(namespace)
	if err != nil {
		return false, err
	}

	if len(sss) != KubePrometheusDefaultStatefulsetCount {
		return false, fmt.Errorf("expect statefulsets count %d, but got count %d now",
			KubePrometheusDefaultStatefulsetCount, len(sss))
	}

	hasNotReadyStatefulset := false
	for _, ss := range sss {
		if kubeClient.IsStatefulsetReady(&ss) {
			log.Printf("the statefulset %s is ready", ss.Name)
			continue
		}
		log.Printf("the statefulset %s is not ready", ss.Name)
		hasNotReadyStatefulset = true
	}

	if hasNotReadyStatefulset {
		return false, fmt.Errorf("some statefulsets are not ready")
	}
	return true, nil
}
