package helminstaller

import (
	"os"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func renderDefaultConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}

	instanceID := helmOptions.InstanceID
	defaultHelmOptions := getDefaultOptionsByInstanceID(instanceID)

	if defaultHelmOptions == nil {
		log.Debugf("Default config for %s wasn't found.", instanceID)
		return options, nil
	}

	log.Infof("Filling default config with instance: %s.", instanceID)
	helmOptions.FillDefaultValue(defaultHelmOptions)
	log.Debugf("Options with default config filled: %v.", helmOptions)

	return types.EncodeStruct(helmOptions)
}

// renderValuesYaml renders options.valuesYaml to options.chart.valuesYaml;
// If options.valuesYaml don't contains ":", it should be a path like "./values.yaml", then read it and transfer to options.chart.valuesYaml
func renderValuesYaml(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}

	// 1. valuesYaml isn't set
	if helmOptions.ValuesYaml == "" {
		return options, nil
	}

	// 2. valuesYaml is a YAML string
	if strings.Contains(helmOptions.ValuesYaml, ": ") {
		helmOptions.Chart.ValuesYaml = helmOptions.ValuesYaml
		helmOptions.ValuesYaml = ""
		return types.EncodeStruct(helmOptions)
	}

	// 3. valuesYaml is a file path
	valuesYamlBytes, err := os.ReadFile(helmOptions.ValuesYaml)
	if err != nil {
		return nil, err
	}
	helmOptions.Chart.ValuesYaml = string(valuesYamlBytes)
	helmOptions.ValuesYaml = ""

	return types.EncodeStruct(helmOptions)
}

func indexStatusGetterFunc(options configmanager.RawOptions) installer.StatusGetterOperation {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		// It's ok to return GetAllResourcesStatus here when err != nil.
		return helm.GetAllResourcesStatus
	}

	instanceID := helmOptions.InstanceID
	statusGetterFunc := getStatusGetterFuncByInstanceID(instanceID)

	if statusGetterFunc == nil {
		return helm.GetAllResourcesStatus
	}

	return statusGetterFunc
}

// e.g. instanceID="argocd-config-001", "argocd" and "argocd-config" both are supported helm charts,
// then DefaultOptionsMap["argocd-config"] needs to be returned.
func getDefaultOptionsByInstanceID(instanceID string) *helm.Options {
	// if string instanceID contains a name, that name is a matched name.
	// e.g. argocd-config-001 contains argocd and argocd-config, so the argocd and argocd-config both are matched name.
	var matchedNames = make([]string, 0)
	for name := range defaults.DefaultOptionsMap {
		if strings.HasPrefix(instanceID, name) {
			matchedNames = append(matchedNames, name)
		}
	}

	if len(matchedNames) == 1 {
		return defaults.DefaultOptionsMap[matchedNames[0]]
	}

	return defaults.DefaultOptionsMap[getLongestMatchedName(matchedNames)]
}

func getLongestMatchedName(nameList []string) string {
	var retStr string
	for _, name := range nameList {
		if len(name) > len(retStr) {
			retStr = name
		}
	}
	return retStr
}

func getStatusGetterFuncByInstanceID(instanceID string) installer.StatusGetterOperation {
	for name, fn := range defaults.StatusGetterFuncMap {
		if strings.Contains(instanceID, name+"-") {
			return fn
		}
	}
	return nil
}
