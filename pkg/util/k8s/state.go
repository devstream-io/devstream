package k8s

import (
	"github.com/devstream-io/devstream/pkg/util/log"
)

type ResourceStatus struct {
	Name  string
	Ready bool
}

type AllResourceStatus struct {
	Deployment  []ResourceStatus
	StatefulSet []ResourceStatus
	DaemonSet   []ResourceStatus
}

//GetResourceStatus get all resource state by input nameSpace and filtermap
func (c *Client) GetResourceStatus(nameSpace string, anFilter map[string]string) (AllResourceStatus, error) {
	stateMap := AllResourceStatus{}
	// 1. list deploy resource
	dps, err := c.ListDeployments(nameSpace)
	if err != nil {
		log.Debugf("Failed to list deployments: %s.", err)
		return stateMap, err
	}

	for _, dp := range dps {
		matchFilterd := filterByAnnotation(dp.GetAnnotations(), anFilter)
		if !matchFilterd {
			log.Debugf("Found unknown deployment: %s.", dp.GetName())
		}
		dpName := dp.GetName()
		ready := c.IsDeploymentReady(&dp)
		stateMap.Deployment = append(stateMap.Deployment, ResourceStatus{dpName, ready})
		log.Debugf("The deployment %s is %t.", dp.GetName(), ready)
	}

	// 2. list statefulsets resource
	sts, err := c.ListStatefulsets(nameSpace)
	if err != nil {
		log.Debugf("Failed to list statefulsets: %s.", err)
		return stateMap, err
	}

	for _, ss := range sts {
		matchFilterd := filterByAnnotation(ss.GetAnnotations(), anFilter)
		if !matchFilterd {
			log.Infof("Found unknown statefulSet: %s.", ss.GetName())
		}

		ready := c.IsStatefulsetReady(&ss)
		ssName := ss.GetName()
		stateMap.StatefulSet = append(stateMap.StatefulSet, ResourceStatus{ssName, ready})
		log.Debugf("The statefulset %s is %t.", ss.GetName(), ready)
	}

	// 3. list daemonset resource
	dss, err := c.ListDaemonsets(nameSpace)
	if err != nil {
		log.Debugf("Failed to list daemonsets: %s.", err)
		return stateMap, err
	}

	for _, ds := range dss {
		matchFilterd := filterByAnnotation(ds.GetAnnotations(), anFilter)
		if matchFilterd {
			log.Infof("Found unknown daemonSet: %s.", ds.GetName())
		}

		ready := c.IsDaemonsetReady(&ds)
		dsName := ds.GetName()
		stateMap.DaemonSet = append(stateMap.DaemonSet, ResourceStatus{dsName, ready})
		log.Debugf("The daemonset %s is %t.", ds.GetName(), ready)
	}
	return stateMap, nil
}

func filterByAnnotation(anInfo map[string]string, anFilter map[string]string) bool {
	for k, v := range anFilter {
		anVal, exist := anInfo[k]
		if !exist || anVal != v {
			return false
		}
	}
	return true
}
