package argocdapp

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

// Reinstall an ArgoCD app
func Reinstall(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	file := defaultYamlPath

	// delete resource
	err = kubectlAction(ActionDelete, file)
	if err != nil {
		return false, err
	}

	// remove app.yaml file
	if err = os.Remove(file); err != nil {
		return false, err
	}

	// recreate app.yaml file
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
