package argocd

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// Install installs ArgoCD with provided options.
func Install(options *map[string]interface{}) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("creating helm client")
	client := createHelmClient(param.Chart.Namespace)
	log.Println("adding and updating argocd helm chart repo")
	addArgoHelmRepo(&client, &param)
	log.Println("installing or updating argocd helm chart")
	installOrUpdateArgoHelmChart(&client, &param)
}
