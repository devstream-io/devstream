package jenkins

import (
	"fmt"
	"regexp"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// validate validates the options provided by the core.
func replaceStroageClass(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}

	// if dev mode, replace the StorageClass name with default StorageClass which is auto created with hostpath type.
	if opts.TestEnv {
		chartNew := opts.Chart
		var err error
		chartNew.ValuesYaml, err = replaceStorageClass(opts.Chart.ValuesYaml)
		if err != nil {
			return nil, err
		}
		opts.Chart = chartNew
	}
	return opts.encode()
}

func replaceStorageClass(valuesYaml string) (string, error) {
	// find the StorageClass name in the options
	re, _ := regexp.Compile(`storageClass:.*\n`)
	storageConfig := re.FindString(valuesYaml)
	if storageConfig == "" {
		return "", fmt.Errorf("storageClass is required in values_yaml config")
	}

	// replace the StorageClass name with default StorageClass name
	valuesYaml = re.ReplaceAllString(valuesYaml, fmt.Sprintf("storageClass: %s\n", jenkinsPvDefaultStorageClassName))
	log.Debugf("new values_yaml whose StorageClass is replaced by default : %s\n", valuesYaml)
	return valuesYaml, nil
}
