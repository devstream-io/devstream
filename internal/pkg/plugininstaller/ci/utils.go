package ci

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/template"
)

// gitlab and jenkins is just a file, so we can just use filename
// github use directory, we shoud process this situation
func getCIFilePath(ciType ciRepoType) string {
	switch ciType {
	case ciGitLabType:
		return ciGitLabConfigLocation
	case ciGitHubType:
		return ciGitHubWorkConfigLocation
	case ciJenkinsType:
		return ciJenkinsConfigLocation
	}
	return ""
}

func filterCIFilesFunc(ciType ciRepoType) file.DirFIleFilterFunc {
	ciFileName := getCIFilePath(ciType)
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		if ciType == ciGitHubType {
			return strings.Contains(filePath, "workflows")
		}
		return filepath.Base(filePath) == ciFileName
	}
}

func processCIFilesFunc(templateName string, vars map[string]interface{}) file.DirFileProcessFunc {
	return func(filePath string) ([]byte, error) {
		if len(vars) == 0 {
			return os.ReadFile(filePath)
		}
		renderContent, err := template.New().FromLocalFile(filePath).SetDefaultRender(templateName, vars).Render()
		if err != nil {
			return nil, err
		}
		return []byte(renderContent), nil
	}
}

func getGitNameFunc(ciType ciRepoType) file.DirFileNameFunc {
	ciFilePath := getCIFilePath(ciType)
	return func(filePath, walkDir string) string {
		fileBaseName := filepath.Base(filePath)
		if ciType == ciGitHubType {
			return filepath.Join(ciFilePath, fileBaseName)
		}
		return ciFilePath
	}
}
