package helminstaller

import (
	"os"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller/defaults"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func renderDefaultConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}

	instanceID := helmOptions.InstanceID

	defaultIns := defaults.GetDefaultHelmAppInstanceByInstanceID(instanceID)
	if defaultIns == nil {
		log.Debugf("Default config for %s wasn't found.", instanceID)
		return options, nil
	}

	log.Infof("Filling default config with instance: %s.", instanceID)
	helmOptions.FillDefaultValue(defaultIns.HelmOptions)
	log.Debugf("Options with default config filled: %v.", helmOptions)

	return mapz.DecodeStructToMap(helmOptions)
}

// renderValuesYaml renders options.valuesYaml to options.chart.valuesYaml;
// If options.valuesYaml doesn't contain ":", it should be a path like "./values.yaml", then read it and transfer to options.chart.valuesYaml
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

	defaultIns := defaults.GetDefaultHelmAppInstanceByInstanceID(helmOptions.InstanceID)

	if defaultIns == nil {
		return helm.GetAllResourcesStatus
	}

	return defaultIns.StatusGetter
}
