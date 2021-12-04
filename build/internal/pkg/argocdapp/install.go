package argocdapp

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// Install creates an ArgoCD app yaml and apply it.
func Install(options *map[string]interface{}) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		log.Fatal(err)
	}

	file := "./app.yaml"
	writeContentToTmpFile(file, appTemplate, &param)
	kubectlApply(file)
}
