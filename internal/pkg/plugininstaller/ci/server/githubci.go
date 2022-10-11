package server

import (
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	ciGitHubType               CIServerType = "github"
	ciGitHubWorkConfigLocation string       = ".github/workflows"
)

type GitHubCI struct {
}

func (g *GitHubCI) Type() CIServerType {
	return ciGitHubType
}

func (g *GitHubCI) CIFilePath(subFilename ...string) string {
	// if subFilename is empty, return dir(.github/workflows)
	if len(subFilename) == 0 {
		return ciGitHubWorkConfigLocation
	}
	// else return dir + subFilename
	return filepath.Join(ciGitHubWorkConfigLocation, filepath.Base(subFilename[0]))
}

func (g *GitHubCI) FilterCIFilesFunc() file.DirFIleFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return strings.Contains(filePath, "workflows")
	}
}

func (g *GitHubCI) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, walkDir string) string {
		return g.CIFilePath(filePath)
	}
}
