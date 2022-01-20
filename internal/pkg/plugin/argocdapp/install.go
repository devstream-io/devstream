package argocdapp

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/util/log"
)

// Install creates an ArgoCD app yaml and apply it.
func Install(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	if errs := validate(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	file := defaultYamlPath
	if err = writeContentToTmpFile(file, appTemplate, &param); err != nil {
		return false, err
	}

	if err = kubectlAction(ActionApply, file); err != nil {
		return false, err
	}

	if err = os.Remove(file); err != nil {
		return false, err
	}

	return true, nil
}
