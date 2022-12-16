package cifile

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

func GetCIFileStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	// init scm client
	client, err := scm.NewClientWithAuth(opts.ProjectRepo)
	if err != nil {
		return nil, err
	}

	// get local file info
	gitMap, err := opts.CIFileConfig.getGitfileMap()
	if err != nil {
		log.Debugf("ci state get gitMap failed: %+v", err)
		return nil, err
	}

	gitFileStatus, err := scm.GetGitFileStats(client, gitMap)
	if err != nil {
		return nil, err
	}
	statusMap := statemanager.ResourceStatus(gitFileStatus)
	return statusMap, nil
}
