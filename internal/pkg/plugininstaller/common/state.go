package common

import (
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func GetPluginAllK8sState(nameSpace string, anFilter, labelFilter map[string]string) (statemanager.ResourceState, error) {
	// 1. init kube client
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	// 2. get all related resource
	allResource, err := kubeClient.GetResourceStatus(nameSpace, anFilter, labelFilter)
	if err != nil {
		log.Debugf("helm status: get status failed: %s", err)
		return nil, err
	}

	// 3. transfer resources status to workflows
	state := &helm.InstanceState{}
	for _, dep := range allResource.Deployment {
		state.Workflows.AddDeployment(dep.Name, dep.Ready)
	}
	for _, sts := range allResource.StatefulSet {
		state.Workflows.AddStatefulset(sts.Name, sts.Ready)
	}
	for _, ds := range allResource.DaemonSet {
		state.Workflows.AddDaemonset(ds.Name, ds.Ready)
	}

	retMap, err := state.ToStringInterfaceMap()
	if err != nil {
		return nil, err
	}
	log.Debugf("Return map: %v.", retMap)
	return retMap, nil
}
