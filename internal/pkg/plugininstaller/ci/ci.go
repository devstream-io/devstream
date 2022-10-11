package ci

import "github.com/devstream-io/devstream/pkg/util/file"

type ciRepoType string

type CI interface {
	// Type return ci type
	Type() ciRepoType
	// CIFilePath returns the file path of ci config file
	// gitlab and jenkins is just a file, so we can just use filename
	// but GitHub use directory, we should process this situation
	// for GitHub: return ".github/workflows" or ".github/workflows/subFilename"
	// for gitlab, jenkins: will ignore subFilename param
	CIFilePath(subFilename ...string) string
	// filterCIFilesFunc returns a filter function to select ci config file
	filterCIFilesFunc() file.DirFIleFilterFunc
	// getGitNameFunc returns a function to transform file path to git name of ci config file
	getGitNameFunc() file.DirFileNameFunc
}

func NewCI(ciType ciRepoType) CI {
	// there are no validation for ciType
	// because we have already validated it by `validate` flag in CIConfig.Type
	switch ciType {
	case ciGitLabType:
		return &GitLabCI{}
	case ciGitHubType:
		return &GitHubCI{}
	case ciJenkinsType:
		return &JenkinsCI{}
	}
	return nil
}
