package argocdapp

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

// Install creates an ArgoCD app yaml and apply it.
func Install(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
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
