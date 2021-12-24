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

	file := "./app.yaml"
	_, errDel := kubectlDelete(file)
	if errDel != nil {
		return false, errDel
	}
	errRemove := os.Remove(file)
	if errRemove != nil {
		return false, errRemove
	}
	return true, nil
}
