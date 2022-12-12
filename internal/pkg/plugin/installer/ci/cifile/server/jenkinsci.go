package server

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/file"
)

const (
	CIJenkinsType           CIServerType = "jenkins"
	ciJenkinsConfigLocation string       = "Jenkinsfile"
)

type JenkinsPipeline struct {
}

// CIFilePath return Jenkinsfile
func (j *JenkinsPipeline) CIFilePath() string {
	return ciJenkinsConfigLocation
}

// FilterCIFilesFunc only get file with name Jenkinsfile
func (j *JenkinsPipeline) FilterCIFilesFunc() file.DirFileFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filepath.Base(filePath) == ciJenkinsConfigLocation
	}
}

func (j *JenkinsPipeline) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, _ string) string {
		return j.CIFilePath()
	}
}
