package helminstaller

import (
	"os"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func RenderDefaultConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}

	instanceID := helmOptions.InstanceID
	defaultHelmOptions := GetDefaultOptionsByInstanceID(instanceID)

	if defaultHelmOptions == nil {
		log.Debugf("Default config for %s wasn't found.", instanceID)
		return options, nil
	}

	log.Infof("Filling default config with instance: %s.", instanceID)
	helmOptions.FillDefaultValue(defaultHelmOptions)
	log.Debugf("Options with default config filled: %v.", helmOptions)

	return types.EncodeStruct(helmOptions)
}

// RenderValuesYaml renders options.valuesYaml to options.chart.valuesYaml;
// If options.valuesYaml don't contains ":", it should be a path like "./values.yaml", then read it and transfer to options.chart.valuesYaml
func RenderValuesYaml(options configmanager.RawOptions) (configmanager.RawOptions, error) {
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

func IndexStatusGetterFunc(options configmanager.RawOptions) plugininstaller.StatusGetterOperation {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		// It's ok to return GetAllResourcesStatus here when err != nil.
		return helm.GetAllResourcesStatus
	}

	instanceID := helmOptions.InstanceID
	statusGetterFunc := GetStatusGetterFuncByInstanceID(instanceID)

	if statusGetterFunc == nil {
		return helm.GetAllResourcesStatus
	}

	return statusGetterFunc
}

func GetDefaultOptionsByInstanceID(instanceID string) *helm.Options {
	for name, options := range defaults.DefaultOptionsMap {
		if strings.Contains(instanceID, name+"-") {
			return options
		}
	}
	return nil
}

func GetStatusGetterFuncByInstanceID(instanceID string) plugininstaller.StatusGetterOperation {
	for name, fn := range defaults.StatusGetterFuncMap {
		if strings.Contains(instanceID, name+"-") {
			return fn
		}
	}
	return nil
}
