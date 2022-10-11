package ci

import "github.com/devstream-io/devstream/pkg/util/file"

const (
	ciJenkinsType           ciRepoType = "jenkins"
	ciJenkinsConfigLocation string     = "Jenkinsfile"
)

type JenkinsCI struct {
}

func (j *JenkinsCI) Type() ciRepoType {
	return ciJenkinsType
}

func (j *JenkinsCI) CIFilePath(_ ...string) string {
	return ciJenkinsConfigLocation
}

func (j *JenkinsCI) filterCIFilesFunc() file.DirFIleFilterFunc {
	return func(filePath string, isDir bool) bool {
		// not process dir
		if isDir {
			return false
		}
		return filePath == ciJenkinsConfigLocation
	}
}

func (j *JenkinsCI) getGitNameFunc() file.DirFileNameFunc {
	return func(filePath, walkDir string) string {
		return j.CIFilePath()
	}
}
