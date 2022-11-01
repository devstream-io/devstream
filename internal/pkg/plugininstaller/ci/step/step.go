package step

import (
	"reflect"

	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
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
	DingTalkSecretKey     string `mapstructure:"DingTalkSecretKey"`
	DingTalkSecretToken   string `mapstructure:"DingTalkSecretToken"`
	ImageRepoSecret       string `mapstructure:"ImageRepoSecret"`
	ImageRepoDockerSecret string `mapstructure:"ImageRepoDockerSecret"`
	CredentialID          string `mapstructure:"StepGlobalVars"`
	SonarqubeSecretKey    string `mapstructure:"SonarqubeSecretKey"`
	GitlabConnectionID    string `mapstructure:"GitlabConnectionID"`
}

// GetStepGlobalVars get global config vars for step
func GetStepGlobalVars(repoInfo *git.RepoInfo) *StepGlobalVars {
	v := &StepGlobalVars{
		ImageRepoSecret:       imageRepoSecretName,
		ImageRepoDockerSecret: imageRepoDockerSecretName,
		DingTalkSecretKey:     dingTalkSecretVal,
		DingTalkSecretToken:   dingTalkSecretToken,
		SonarqubeSecretKey:    sonarSecretKey,
		GitlabConnectionID:    gitlabConnectionName,
	}
	if repoInfo.IsGitlab() && repoInfo.SSHPrivateKey != "" {
		v.CredentialID = gitlabCredentialName
	} else if repoInfo.IsGithub() {
		v.CredentialID = githubCredentialName
	}
	return v
}

func ExtractValidStepConfig(p any) []StepConfigAPI {
	var stepConfigs []StepConfigAPI
	// 1. add pipeline plugin
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
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

func GenerateCIFileVars(p any, repoInfo *git.RepoInfo) cifile.CIFileVarsMap {
	varMap, _ := mapz.DecodeStructToMap(p)
	globalVarsMap, _ := mapz.DecodeStructToMap(
		GetStepGlobalVars(repoInfo),
	)
	err := mergo.Merge(&varMap, globalVarsMap)
	if err != nil {
		log.Warnf("cifile merge CIFileVarsMap failed: %+v", err)
	}
	return varMap
}
