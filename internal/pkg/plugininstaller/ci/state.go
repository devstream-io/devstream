package ci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

func GetCIFileStatus(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	fileLocation := getCIFilePath(opts.CIConfig.Type)
	client, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return nil, err
	}
	gitFileInfo, err := client.GetLocationInfo(fileLocation)
	if err != nil {
		return nil, err
	}
	statusMap := make(map[string]interface{})
	for _, item := range gitFileInfo {
		statusMap[item.Path] = map[string]string{
			"sha":    item.SHA,
			"path":   item.Path,
			"branch": item.Branch,
		}
	}
	return statusMap, nil
}
