package server

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	ciGitLabType           CIServerType = "gitlab"
	ciGitLabConfigLocation string       = ".gitlab-ci.yml"
)

type GitLabCI struct {
}

func (g *GitLabCI) Type() CIServerType {
	return ciGitLabType
}

func (g *GitLabCI) CIFilePath(_ ...string) string {
	return ciGitLabConfigLocation
}

func (g *GitLabCI) FilterCIFilesFunc() file.DirFIleFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filepath.Base(filePath) == g.CIFilePath()
	}
}

func (g *GitLabCI) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, walkDir string) string {
		return g.CIFilePath()
	}
}
