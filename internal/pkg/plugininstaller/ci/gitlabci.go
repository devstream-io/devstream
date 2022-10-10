package ci

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	ciGitLabType           ciRepoType = "gitlab"
	ciGitLabConfigLocation string     = ".gitlab-ci.yml"
)

type GitLabCI struct {
}

func (g *GitLabCI) Type() ciRepoType {
	return ciGitLabType
}

func (g *GitLabCI) CIFilePath(_ ...string) string {
	return ciGitLabConfigLocation
}

func (g *GitLabCI) filterCIFilesFunc() file.DirFIleFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filepath.Base(filePath) == g.CIFilePath()
	}
}

func (g *GitLabCI) getGitNameFunc() file.DirFileNameFunc {
	return func(filePath, walkDir string) string {
		return g.CIFilePath()
	}
}
