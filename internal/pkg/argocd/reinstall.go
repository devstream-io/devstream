package argocd

import (
	"log"
)

func Reinstall(options *map[string]interface{}) (bool, error) {
	acd, err := NewArgoCD(options)
	if err != nil {
		return false, err
	}

	log.Println("uninstalling argocd helm chart")
	if err := acd.uninstallHelmChart(); err != nil {
		return false, err
	}

	log.Println("installing argocd helm chart")
	if err := acd.installOrUpgradeHelmChart(); err != nil {
		return false, err
	}

	return true, nil
}
