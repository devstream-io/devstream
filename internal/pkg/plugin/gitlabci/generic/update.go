package generic

import "github.com/devstream-io/devstream/internal/pkg/statemanager"

func Update(options map[string]interface{}) (statemanager.ResourceStatus, error) {
	return Create(options)
}
