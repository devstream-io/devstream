package step

import (
	"reflect"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

type StepConfigAPI interface {
	GetJenkinsPlugins() []*jenkins.JenkinsPlugin
	ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error)
	ConfigGithub(client *github.Client) error
}

type StepGlobalOption struct {
	RepoInfo *git.RepoInfo
}

type StepGlobalVars struct {
	ImageRepoSecret   string `mapstructure:"ImageRepoSecret"`
	DingTalkSecretKey string `mapstructure:"DingTalkSecretKey"`
	CredentialID      string `mapstructure:"StepGlobalVars"`
}

// GetStepGlobalVars get global config vars for step
func GetStepGlobalVars(repoInfo *git.RepoInfo) *StepGlobalVars {
	v := &StepGlobalVars{
		ImageRepoSecret:   imageRepoSecretName,
		DingTalkSecretKey: dingTalkSecretKey,
	}
	if repoInfo.IsGitlab() && repoInfo.SSHPrivateKey != "" {
		v.CredentialID = gitlabCredentialName
	} else if repoInfo.IsGithub() {
		v.CredentialID = sshKeyCredentialName
	}
	return v
}

func ExtractValidStepConfig(p any) []StepConfigAPI {
	var stepConfigs []StepConfigAPI
	// 1. add pipeline plugin
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			fieldVal, ok := valueField.Interface().(StepConfigAPI)
			if ok {
				stepConfigs = append(stepConfigs, fieldVal)
			} else {
				log.Warnf("jenkins extract pipeline plugins failed: %+v", valueField)
			}
		}
	}
	return stepConfigs
}

func GetRepoStepConfig(repoInfo *git.RepoInfo) []StepConfigAPI {
	var stepConfigs []StepConfigAPI
	plugGlobalConfig := &StepGlobalOption{
		RepoInfo: repoInfo,
	}
	switch repoInfo.RepoType {
	case "gitlab":
		stepConfigs = append(stepConfigs, newGitlabStep(plugGlobalConfig))
	case "github":
		stepConfigs = append(stepConfigs, newGithubStep(plugGlobalConfig))
	}
	return stepConfigs
}
