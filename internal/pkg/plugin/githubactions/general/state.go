package general

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func getState(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// get ci sate
	return cifile.GetCIFileStatus(options)
}
