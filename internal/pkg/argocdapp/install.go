package argocdapp

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// Install creates an ArgoCD app yaml and apply it.
func Install(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		log.Fatal(err)
	}

	file := "./app.yaml"
	writeContentToTmpFile(file, appTemplate, &param)
	_, errApply := kubectlApply(file)
	if errApply != nil {
		return false, errApply
	}

	return true, nil
}
