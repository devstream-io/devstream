package argocdapp

import (
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
