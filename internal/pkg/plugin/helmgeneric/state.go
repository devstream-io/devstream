package helmgeneric

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

// return empty
func getEmptyStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	retStatus := make(statemanager.ResourceStatus)
	return retStatus, nil
}
