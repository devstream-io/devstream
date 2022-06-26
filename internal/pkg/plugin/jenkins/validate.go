package jenkins

import (
	"fmt"
	"regexp"

	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// validate validates the options provided by the core.
func validate(opts *Options) []error {
	errs := helm.Validate(opts.GetHelmParam())
	if len(errs) != 0 {
		return errs
	}

	// if dev mode, replace the storage class name with default storage class which is auto created with hostpath type.
	if opts.TestEnv {
		chartNew := opts.Chart
		var err error
		chartNew.ValuesYaml, err = ReplaceStorageClass(opts.Chart.ValuesYaml)
		if err != nil {
			return []error{err}
		}
		opts.Chart = chartNew
	}
	return nil
}

func ReplaceStorageClass(valuesYaml string) (string, error) {
	// find the storage class name in the options
	re, _ := regexp.Compile(`storageClass:.*\n`)
	storageConfig := re.FindString(valuesYaml)
	if storageConfig == "" {
		return "", fmt.Errorf("storageClass is required in  values_yaml config")
	}

	// replace the storage class name with default storage class name
	valuesYaml = re.ReplaceAllString(valuesYaml, fmt.Sprintf("storageClass: %s\n", JenkinsPvDefaultStorageClassName))
	log.Debugf("new values_yaml whose storage class is replaced by default : %s\n", valuesYaml)

	return valuesYaml, nil
}
