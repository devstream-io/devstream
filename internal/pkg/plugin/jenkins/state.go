package jenkins

import (
	"fmt"

	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/log"
)

var DefaultStatefulsetTplList = []string{
	// ${release-name}-jenkins
	"%s-jenkins",
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
	DefaultStatefulsetList := GetDefaultStatefulsetList(releaseName)

	for _, ssName := range DefaultStatefulsetList {
		retState.Workflows.AddStatefulset(ssName, true)
	}

	return retState
}
