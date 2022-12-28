package server

import (
	"github.com/devstream-io/devstream/pkg/util/file"
)

type CIServerType string

type CIServerOptions interface {
	// CIFilePath returns the file path of ci config file
	// gitlab and jenkins is just a file, so we can just use filename
	// but GitHub use directory, we should process this situation
	// for GitHub: return ".github/workflows" or ".github/workflows/subFilename"
	// for gitlab, jenkins: will ignore subFilename param
	CIFilePath() string

	// FilterCIFilesFunc returns a filter function to select ci config file
	FilterCIFilesFunc() file.DirFileFilterFunc
	// GetGitNameFunc returns a function to transform file path to git name of ci config file
	GetGitNameFunc() file.DirFileNameFunc
}

func NewCIServer(ciType CIServerType) CIServerOptions {
	// there are no validation for ciType
	// because we have already validated it by `validate` flag in CIFileConfig.Type
	switch ciType {
	case CIGitLabType:
		return &GitLabCI{}
	case CIGithubType:
		return &GitHubWorkflow{}
	case CIJenkinsType:
		return &JenkinsPipeline{}
	}
	//TODO(jiafeng meng): if ciType is not exist, this will cause error
	return nil
}
