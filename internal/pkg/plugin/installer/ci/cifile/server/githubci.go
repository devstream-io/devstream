package server

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	CIGithubType               CIServerType = "github"
	ciGitHubWorkConfigLocation string       = ".github/workflows"
	ciGithubTempName           string       = "app.yaml"
)

type GitHubWorkflow struct {
}

// CIFilePath return .github/workflows/app.yml
func (g *GitHubWorkflow) CIFilePath() string {
	return filepath.Join(ciGitHubWorkConfigLocation, ciGithubTempName)
}

func (g *GitHubWorkflow) FilterCIFilesFunc() file.DirFileFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return strings.Contains(filePath, "yml") || strings.Contains(filePath, "yaml")
	}
}

func (g *GitHubWorkflow) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, _ string) string {
		basePath := filepath.Base(filePath)
		return path.Join(ciGitHubWorkConfigLocation, basePath)
	}
}
