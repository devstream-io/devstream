package argocd

import (
	"log"
)

// Install installs ArgoCD with provided options.
func Install(options *map[string]interface{}) (bool, error) {
	acd, err := NewArgoCD(options)
	if err != nil {
		return false, err
	}

	log.Println("Installing or updating argocd helm chart ...")
	if err := acd.installOrUpgradeHelmChart(); err != nil {
		return false, err
	}

	return true, nil
}
