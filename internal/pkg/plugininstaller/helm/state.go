package helm

import (
	"fmt"

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
			dpAn := dp.GetAnnotations()
			helmNameAn, exist := dpAn[helm.GetAnnotationName()]
			if !exist || helmNameAn != opts.GetReleaseName() {
				log.Infof("Found unknown deployment: %s.", dp.GetName())
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

// GetPlugAllStateWrapper will get deploy, ds, statefulset status
func GetPluginAllState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
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
		log.Debugf("Failed to list deployments: %s.", err)
		return nil, err
	}

	state := &helm.InstanceState{}
	for _, dp := range dps {
		dpName := dp.GetName()
		dpAn := dp.GetAnnotations()
		helmNameAn, exist := dpAn[helm.GetAnnotationName()]
		if !exist || helmNameAn != opts.GetReleaseName() {
			log.Infof("Found unknown deployment: %s.", dp.GetName())
		}

		ready := kubeClient.IsDeploymentReady(&dp)
		state.Workflows.AddDeployment(dpName, ready)
		log.Debugf("The deployment %s is %t.", dp.GetName(), ready)
	}

	sts, err := kubeClient.ListStatefulsets(opts.GetNamespace())
	if err != nil {
		log.Debugf("Failed to list statefulsets: %s.", err)
		return nil, err
	}

	for _, ss := range sts {
		ssName := ss.GetName()
		ssAn := ss.GetAnnotations()
		helmNameAn, exist := ssAn[helm.GetAnnotationName()]
		if !exist || helmNameAn != opts.GetReleaseName() {
			log.Infof("Found unknown statefulSet: %s.", ss.GetName())
		}

		ready := kubeClient.IsStatefulsetReady(&ss)
		state.Workflows.AddStatefulset(ssName, ready)
		log.Debugf("The statefulset %s is %t.", ss.GetName(), ready)
	}

	dss, err := kubeClient.ListDaemonsets(opts.GetNamespace())
	if err != nil {
		log.Debugf("Failed to list daemonsets: %s.", err)
		return nil, err
	}

	for _, ds := range dss {
		dsName := ds.GetName()
		ssAn := ds.GetAnnotations()
		helmNameAn, exist := ssAn[helm.GetAnnotationName()]
		if !exist || helmNameAn != opts.GetReleaseName() {
			log.Infof("Found unknown daemonSet: %s.", ds.GetName())
		}

		ready := kubeClient.IsDaemonsetReady(&ds)
		state.Workflows.AddDaemonset(dsName, ready)
		log.Debugf("The daemonset %s is %t.", ds.GetName(), ready)
	}

	retMap := state.ToStringInterfaceMap()
	log.Debugf("Return map: %v.", retMap)
	return retMap, nil
}

//  GetPluginStaticStateByReleaseNameWrapper will return status by static deploy config
func GetPluginStaticStateByReleaseNameWrapper(defaultDepList []string) plugininstaller.StatusOperation {
	getPluginStaticState := func(options plugininstaller.RawOptions) (map[string]interface{}, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}
		retState := &helm.InstanceState{}
		for _, dpName := range defaultDepList {
			namedDpName := fmt.Sprintf(dpName, opts.GetReleaseName())
			retState.Workflows.AddDeployment(namedDpName, true)
		}
		return retState.ToStringInterfaceMap(), nil
	}
	return getPluginStaticState
}
