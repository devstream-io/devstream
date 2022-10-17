package helmgeneric

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

// return empty
func getEmptyStatus(options plugininstaller.RawOptions) (statemanager.ResourceStatus, error) {
	retStatus := make(statemanager.ResourceStatus)
	return retStatus, nil
}
