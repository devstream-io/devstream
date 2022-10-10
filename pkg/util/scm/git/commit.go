package git

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type GitFileContentMap map[string][]byte

type CommitInfo struct {
	CommitMsg    string
	CommitBranch string
	GitFileMap   GitFileContentMap
}

// GitFilePathInfo contains file local path and remote git path
type GitFilePathInfo struct {
	SourcePath      string
	DestinationPath string
}

// unused
func GetFileContent(files []*GitFilePathInfo) GitFileContentMap {
	gitFileMap := make(map[string][]byte)
	for _, f := range files {
		content, err := os.ReadFile(f.SourcePath)
		if err != nil {
			log.Warnf("Repo Process file content error: %s", err)
			continue
		}
		if f.DestinationPath == "" {
			log.Warnf("Repo file destination path is not set")
			continue
		}
		gitFileMap[f.DestinationPath] = content
	}
	return gitFileMap
}

// unused
func GenerateGitFileInfo(filePaths []string, gitDirPath string) ([]*GitFilePathInfo, error) {
	gitFileInfos := make([]*GitFilePathInfo, 0)
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			return gitFileInfos, err
		}
		if !info.IsDir() {
			gitFileInfos = append(gitFileInfos, &GitFilePathInfo{
				SourcePath:      filePath,
				DestinationPath: filepath.Join(gitDirPath, filePath),
			})
			continue
		}
		err = filepath.Walk(filePath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Debugf("git files walk error: %s.", err)
				return err
			}
			// not process dir
			if info.IsDir() {
				return nil
			}
			repoPath, err := filepath.Rel(filePath, path)
			if err != nil {
				log.Debugf("git files get relative path error: %s", err)
			}
			if gitDirPath != "" {
				repoPath = filepath.Join(gitDirPath, repoPath)
			}
			gitFileInfos = append(gitFileInfos, &GitFilePathInfo{
				SourcePath:      path,
				DestinationPath: repoPath,
			})
			return nil
		})
		if err != nil {
			return gitFileInfos, err
		}
	}

	return gitFileInfos, nil
}
