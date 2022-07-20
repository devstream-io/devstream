package helm

import (
	"k8s.io/utils/strings/slices"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// GetPluginStaticStateWrapper will return status by static deploy config
func GetPluginStaticStateWrapper(defaultDepList []string) plugininstaller.StatusOperation {
	getPluginStaticState := func(options plugininstaller.RawOptions) (map[string]interface{}, error) {
		retState := &helm.InstanceState{}
		for _, dpName := range defaultDepList {
			retState.Workflows.AddDeployment(dpName, true)
		}
		return retState.ToStringInterfaceMap(), nil
	}
	return getPluginStaticState

}

// GetPluginDynamicStateWrapper will return status by query kubernetes deployment
func GetPluginDynamicStateWrapper(defaultDepList []string) plugininstaller.StatusOperation {
	getPluginDynamicState := func(options plugininstaller.RawOptions) (map[string]interface{}, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}

		kubeClient, err := k8s.NewClient()
		if err != nil {
			return nil, err
		}

		dps, err := kubeClient.ListDeployments(opts.GetNamespace())
		if err != nil {
			return nil, err
		}

		retState := &helm.InstanceState{}
		for _, dp := range dps {
			dpName := dp.GetName()
			if !slices.Contains(defaultDepList, dpName) {
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
	return getPluginDynamicState
}
