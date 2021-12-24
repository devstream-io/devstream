package argocdapp

import (
	"os"

	"github.com/mitchellh/mapstructure"
)

// Reinstall an ArgoCD app .
func Reinstall(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	file := "./app.yaml"

	//delete resource
	_, errDel := kubectlDelete(file)
	if errDel != nil {
		return false, errDel
	}

	//remove app.yaml file
	errRemove := os.Remove(file)
	if errRemove != nil {
		return false, err
	}

	//recreate  app.yaml file
	writeContentToTmpFile(file, appTemplate, &param)
	_, errApply := kubectlApply(file)
	if errApply != nil {
		return false, errApply
	}

	return true, nil
}
