package ci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func GetCIFileStatus(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	fileLocation := getCIFilePath(opts.CIConfig.Type)
	client, err := opts.ProjectRepo.NewClient()
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
