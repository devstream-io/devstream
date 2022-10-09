package helm

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// GetPluginAllState will get the State of k8s Deployment, DaemonSet and StatefulSet resources
func GetPluginAllState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		log.Debugf("helm init k8s client to get state failed: %+v", err)
		return nil, err
	}

	anFilter := map[string]string{
		helm.GetAnnotationName(): opts.GetReleaseName(),
	}
	labelFilter := map[string]string{}
	allResourceState, err := kubeClient.GetResourceStatus(opts.GetNamespace(), anFilter, labelFilter)
	if err != nil {
		log.Debugf("helm get resource state failed: %+v", err)
		return nil, err
	}
	return allResourceState.ToStringInterfaceMap()
}
