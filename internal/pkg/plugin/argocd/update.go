package argocd

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	_, err := Delete(options)
	if err != nil {
		log.Errorf("Failed to delete the ArgoCD: %s.", err)
		return nil, err
	}
	return Create(options)
}
