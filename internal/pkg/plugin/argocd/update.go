package argocd

import (
	"time"

	"github.com/merico-dev/stream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
	_, err := Delete(options)
	if err != nil {
		log.Errorf("Failed to delete the ArgoCD: %s.", err)
		return nil, err
	}

	<-time.NewTicker(3 * time.Second).C
	return Create(options)
}
