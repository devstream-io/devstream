package helm

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// GetPlugAllStateWrapper will get deploy, ds, statefulset status
func GetPluginAllState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	// 1. init kube client
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	anFilter := map[string]string{
		helm.GetAnnotationName(): opts.GetReleaseName(),
	}

	// 2. get all related resource
	allResource, err := kubeClient.GetResourceStatus(opts.GetNamespace(), anFilter)
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
		state.Workflows.AddStatefulset(ds.Name, ds.Ready)
	}

	retMap := state.ToStringInterfaceMap()
	log.Debugf("Return map: %v.", retMap)
	return retMap, nil
}
