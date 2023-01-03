package defaults

import (
	"strings"
	"sync"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
)

type HelmAppInstance struct {
	Name         string
	HelmOptions  *helm.Options
	StatusGetter installer.StatusGetterOperation
}

var (
	helmAppInstances map[string]*HelmAppInstance
	once             sync.Once
)

func RegisterDefaultHelmAppInstance(name string, options *helm.Options, statusGetter installer.StatusGetterOperation) {
	// make sure the map is initialized only once
	once.Do(func() {
		helmAppInstances = make(map[string]*HelmAppInstance)
	})

	helmAppInstances[name] = &HelmAppInstance{
		Name:         name,
		HelmOptions:  options,
		StatusGetter: statusGetter,
	}
}

// GetDefaultHelmAppInstanceByInstanceID will return the default helm app instance
// by matching the prefix of given instanceID and the list of helm app instances names.
// It will return the longest matched name if there are multiple matched names, and return nil if no matched name is found.
// e.g. instanceID="argocd-config-001", "argocd" and "argocd-config" both are supported helm charts,
// then the HelmAppInstance named "argocd-config" will be returned.
func GetDefaultHelmAppInstanceByInstanceID(instanceID string) *HelmAppInstance {
	longestMatchedName := ""
	for curNameToMatch := range helmAppInstances {
		if strings.HasPrefix(instanceID, curNameToMatch) && len(curNameToMatch) > len(longestMatchedName) {
			longestMatchedName = curNameToMatch
		}
	}

	if longestMatchedName != "" {
		return helmAppInstances[longestMatchedName]
	}
	return nil
}
