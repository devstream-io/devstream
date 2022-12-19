package scm

import (
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func GetGitFileStats(client ClientOperation, gitFiles git.GitFileContentMap) (map[string]any, error) {
	status := make(map[string]any)
	for scmPath, content := range gitFiles {
		// get localSHA
		localFileSHA := calculateSHA(content)
		// get scmInfo
		gitFileInfos, err := client.GetPathInfo(scmPath)
		if err != nil {
			log.Debugf("ci status get location info failed: %+v", err)
			return nil, err
		}
		scmFileStatus := []map[string]string{}
		for _, fileStatus := range gitFileInfos {
			scmFileStatus = append(scmFileStatus, map[string]string{
				"scmSHA":    fileStatus.SHA,
				"scmBranch": fileStatus.Branch,
			})
		}
		status[scmPath] = map[string]any{
			"localSHA": localFileSHA,
			"scm":      scmFileStatus,
		}
	}
	return status, nil
}
