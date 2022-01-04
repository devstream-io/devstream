package argocdapp

import (
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
	err = writeContentToTmpFile(file, appTemplate, &param)
	if err != nil {
		return false, err
	}

	err = kubectlAction(ActionApply, file)
	if err != nil {
		return false, err
	}

	return true, nil
}
