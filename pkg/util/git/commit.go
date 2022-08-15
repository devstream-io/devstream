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

// GitFileInfo contains file local path and remote git path
type GitFilePathInfo struct {
	SourcePath       string
	DestionationPath string
}

func GetFileContent(files []*GitFilePathInfo) GitFileContentMap {
	gitFileMap := make(map[string][]byte)
	for _, f := range files {
		content, err := os.ReadFile(f.SourcePath)
		if err != nil {
			log.Warnf("Repo Process file content error: %s", err)
			continue
		}
		if f.DestionationPath == "" {
			log.Warnf("Repo file destination path is not set")
			continue
		}
		gitFileMap[f.DestionationPath] = content
	}
	return gitFileMap
}

func GenerateGitFileInfo(filePaths []string, gitDirPath string) ([]*GitFilePathInfo, error) {
	gitFileInfos := make([]*GitFilePathInfo, 0)
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			return gitFileInfos, err
		}
		if !info.IsDir() {
			gitFileInfos = append(gitFileInfos, &GitFilePathInfo{
				SourcePath:       filePath,
				DestionationPath: filepath.Join(gitDirPath, filePath),
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
				SourcePath:       path,
				DestionationPath: repoPath,
			})
			return nil
		})
		if err != nil {
			return gitFileInfos, err
		}
	}

	return gitFileInfos, nil
}
