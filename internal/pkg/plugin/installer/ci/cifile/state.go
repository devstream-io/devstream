package cifile

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
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

	statusMap := make(statemanager.ResourceStatus)
	for scmPath, content := range gitMap {
		localFileSHA := git.CalculateLocalFileSHA(content)
		// get remote file status
		statusMap[scmPath] = map[string]interface{}{
			"localSHA": localFileSHA,
			"scm":      getSCMFileStatus(client, scmPath),
		}
	}
	return statusMap, nil
}

func getSCMFileStatus(client scm.ClientOperation, scmPath string) (scmFileStatus []map[string]string) {
	gitFileInfos, err := client.GetPathInfo(scmPath)
	if err != nil {
		log.Debugf("ci status get location info failed: %+v", err)
		return scmFileStatus
	}
	for _, fileStatus := range gitFileInfos {
		scmFileStatus = append(scmFileStatus, map[string]string{
			"scmSHA":    fileStatus.SHA,
			"scmBranch": fileStatus.Branch,
		})
	}
	return scmFileStatus
}
