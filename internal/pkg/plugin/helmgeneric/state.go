package helmgeneric

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

// return empty
func getEmptyState(options plugininstaller.RawOptions) (statemanager.ResourceStatus, error) {
	retMap := make(statemanager.ResourceStatus)
	return retMap, nil
}
