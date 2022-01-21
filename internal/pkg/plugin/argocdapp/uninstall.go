package argocdapp

import (
	"fmt"
	"github.com/merico-dev/stream/internal/pkg/log"
	"os"

	"github.com/mitchellh/mapstructure"
)

// Uninstall uninstall an ArgoCD app by yaml.
func Uninstall(options *map[string]interface{}) (bool, error) {
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

	if err = kubectlAction(ActionDelete, file); err != nil {
		return false, err
	}

	if err = os.Remove(file); err != nil {
		return false, err
	}

	return true, nil
}
