package server

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	CIGitLabType           CIServerType = "gitlab"
	ciGitLabConfigLocation string       = ".gitlab-ci.yml"
)

type GitLabCI struct {
}

// CIFilePath return .gitlab-ci.yml
func (g *GitLabCI) CIFilePath() string {
	return ciGitLabConfigLocation
}

func (g *GitLabCI) FilterCIFilesFunc() file.DirFileFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filepath.Base(filePath) == ciGitLabConfigLocation
	}
}

func (g *GitLabCI) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, _ string) string {
		return g.CIFilePath()
	}
}
