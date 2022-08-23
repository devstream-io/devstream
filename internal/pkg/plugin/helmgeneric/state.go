package helmgeneric

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

// return empty
func getEmptyState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	retMap := make(statemanager.ResourceState)
	return retMap, nil
}
