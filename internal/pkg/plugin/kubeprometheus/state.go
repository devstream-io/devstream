package kubeprometheus

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func getDefaultDeploymentList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range defaultDeploymentTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultDeploymentList: %v", retList)

	return retList
}

func getDefaultDaemonsetList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range defaultDaemonsetTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultDaemonsetList: %v", retList)

	return retList
}

func getDefaultStatefulsetList(releaseName string) []string {
	retList := make([]string, 0)

	for _, name := range defaultStatefulsetTplList {
		retList = append(retList, fmt.Sprintf(name, releaseName))
	}
	log.Debugf("DefaultStatefulsetList: %v", retList)

	return retList
}

func getStaticState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	retState := &helmCommon.InstanceState{}
	DefaultDeploymentList := getDefaultDeploymentList(opts.GetReleaseName())
	DefaultDaemonsetList := getDefaultDaemonsetList(opts.GetReleaseName())
	DefaultStatefulsetList := getDefaultStatefulsetList(opts.GetReleaseName())

	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	for _, dsName := range DefaultDaemonsetList {
		retState.Workflows.AddDaemonset(dsName, true)
	}
	for _, ssName := range DefaultStatefulsetList {
		retState.Workflows.AddStatefulset(ssName, true)
	}

	return retState.ToStringInterfaceMap(), nil
}
