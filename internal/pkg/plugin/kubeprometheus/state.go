package kubeprometheus

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var DefaultDeploymentTplList = []string{
	// ${release-name}-grafana
	"%s-grafana",
	// ${release-name}-kube-prometheus-stack-operator
	"%s-kube-prometheus-stack-operator",
	// ${release-name}-kube-state-metrics
	"%s-kube-state-metrics",
}

var DefaultDaemonsetTplList = []string{
	// ${release-name}-prometheus-node-exporter
	"%s-prometheus-node-exporter",
}

var DefaultStatefulsetTplList = []string{
	// alertmanager-${release-name}-kube-prometheus-stack-alertmanager
	"alertmanager-%s-kube-prometheus-stack-alertmanager",
	// prometheus-${release-name}-kube-prometheus-stack-prometheus
	"prometheus-%s-kube-prometheus-stack-prometheus",
}

func GetDefaultDeploymentList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range DefaultDeploymentTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultDeploymentList: %v", retList)

	return retList
}

func GetDefaultDaemonsetList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range DefaultDaemonsetTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultDaemonsetList: %v", retList)

	return retList
}

func GetDefaultStatefulsetList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range DefaultStatefulsetTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultStatefulsetList: %v", retList)

	return retList
}

func GetStaticState(releaseName string) *helm.InstanceState {
	retState := &helm.InstanceState{}
	DefaultDeploymentList := GetDefaultDeploymentList(releaseName)
	DefaultDaemonsetList := GetDefaultDaemonsetList(releaseName)
	DefaultStatefulsetList := GetDefaultStatefulsetList(releaseName)

	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	for _, dsName := range DefaultDaemonsetList {
		retState.Workflows.AddDaemonset(dsName, true)
	}
	for _, ssName := range DefaultStatefulsetList {
		retState.Workflows.AddStatefulset(ssName, true)
	}

	return retState
}
