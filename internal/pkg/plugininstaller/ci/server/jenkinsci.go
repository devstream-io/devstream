package server

import "github.com/devstream-io/devstream/pkg/util/file"

const (
	ciJenkinsType           CIServerType = "jenkins"
	ciJenkinsConfigLocation string       = "Jenkinsfile"
)

type JenkinsCI struct {
}

func (j *JenkinsCI) Type() CIServerType {
	return ciJenkinsType
}

func (j *JenkinsCI) CIFilePath(_ ...string) string {
	return ciJenkinsConfigLocation
}

func (j *JenkinsCI) FilterCIFilesFunc() file.DirFIleFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filePath == ciJenkinsConfigLocation
	}
}

func (j *JenkinsCI) GetGitNameFunc() file.DirFileNameFunc {
	return func(filePath, walkDir string) string {
		return j.CIFilePath()
	}
}
