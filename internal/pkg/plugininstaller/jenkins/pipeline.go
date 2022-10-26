package jenkins

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type pipeline struct {
	Job             string                          `mapstructure:"jobName" validate:"required"`
	JenkinsfilePath string                          `mapstructure:"jenkinsfilePath" validate:"required"`
	ImageRepo       *plugins.ImageRepoJenkinsConfig `mapstructure:"imageRepo"`
	Dingtalk        *plugins.DingtalkJenkinsConfig  `mapstructure:"dingTalk"`
	Sonarqube       *plugins.SonarQubeJenkinsConfig `mapstructure:"sonarqube"`
	Custom          map[string]interface{}          `mapstructure:",remain"`
}

func (p *pipeline) getJobName() string {
	if strings.Contains(p.Job, "/") {
		return strings.Split(p.Job, "/")[1]
	}
	return p.Job
}

func (p *pipeline) getJobFolder() string {
	if strings.Contains(p.Job, "/") {
		return strings.Split(p.Job, "/")[0]
	}
	return ""
}

func (p *pipeline) extractPlugins(repoInfo *git.RepoInfo) []plugins.PluginConfigAPI {
	var pluginConfigs []plugins.PluginConfigAPI
	plugGlobalConfig := &plugins.PluginGlobalConfig{
		RepoInfo: repoInfo,
	}
	// 1. add pipeline plugin
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		if valueField.Kind() == reflect.Ptr && !valueField.IsNil() {
			fieldVal, ok := valueField.Interface().(plugins.PluginConfigAPI)
			if ok {
				pluginConfigs = append(pluginConfigs, fieldVal)
			} else {
				log.Warnf("jenkins extract pipeline plugins failed: %+v", valueField)
			}
		}

	}
	// 2. add scm plugin
	switch repoInfo.RepoType {
	case "gitlab":
		pluginConfigs = append(pluginConfigs, plugins.NewGitlabPlugin(plugGlobalConfig))
	case "github":
		pluginConfigs = append(pluginConfigs, plugins.NewGithubPlugin(plugGlobalConfig))
	}
	return pluginConfigs
}

func (p *pipeline) install(jenkinsClient jenkins.JenkinsAPI, repoInfo *git.RepoInfo, secretToken string) error {
	// 1. install jenkins plugins
	pipelinePlugins := p.extractPlugins(repoInfo)
	if err := plugins.EnsurePluginInstalled(jenkinsClient, pipelinePlugins); err != nil {
		return err
	}
	// 2. config plugins by casc
	if err := plugins.ConfigPlugins(jenkinsClient, pipelinePlugins); err != nil {
		return err
	}
	// 3. create or update jenkins job
	return p.createOrUpdateJob(jenkinsClient, repoInfo, secretToken)
}

func (p *pipeline) remove(jenkinsClient jenkins.JenkinsAPI, repoInfo *git.RepoInfo) error {
	job, err := jenkinsClient.GetFolderJob(p.getJobName(), p.getJobFolder())
	if err != nil {
		// check job is already been deleted
		if jenkins.IsNotFoundError(err) {
			return nil
		}
		return err
	}
	_, err = job.Delete(context.Background())
	log.Debugf("jenkins delete job %s status: %v", p.Job, err)
	return nil
}

func (p *pipeline) createOrUpdateJob(jenkinsClient jenkins.JenkinsAPI, repoInfo *git.RepoInfo, secretToken string) error {
	// 1. render groovy script
	jobRenderInfo := &jenkins.JobScriptRenderInfo{
		RepoType:          repoInfo.RepoType,
		JobName:           p.getJobName(),
		RepositoryURL:     repoInfo.CloneURL,
		Branch:            repoInfo.Branch,
		SecretToken:       secretToken,
		FolderName:        p.getJobFolder(),
		GitlabConnection:  plugins.GitlabConnectionName,
		RepoURL:           repoInfo.BuildScmURL(),
		RepoOwner:         repoInfo.GetRepoOwner(),
		RepoName:          repoInfo.Repo,
		RepoCredentialsId: plugins.GetRepoCredentialsId(repoInfo),
	}
	jobScript, err := jenkins.BuildRenderedScript(jobRenderInfo)
	if err != nil {
		log.Debugf("jenkins render template failed: %s", err)
		return err
	}
	// 2. execute script to create jenkins job
	_, err = jenkinsClient.ExecuteScript(jobScript)
	if err != nil {
		log.Debugf("jenkins execute script failed: %s", err)
		return err
	}
	return nil
}

func (p *pipeline) checkValid() error {
	if strings.Contains(p.Job, "/") {
		strs := strings.Split(p.Job, "/")
		if len(strs) != 2 || len(strs[0]) == 0 || len(strs[1]) == 0 {
			return fmt.Errorf("jenkins jobName illegal: %s", p.Job)
		}
	}
	return nil
}

func (p *pipeline) buildCIConfig(repoInfo *git.RepoInfo) (*ci.CIConfig, error) {
	ciConfig := &ci.CIConfig{
		Type:           "jenkins",
		ConfigLocation: p.JenkinsfilePath,
	}
	// update ci render variables by plugins
	pipelinePlugins := p.extractPlugins(repoInfo)
	jenkinsRenderVars := plugins.GetPluginsRenderVariables(pipelinePlugins)
	jenkinsRenderVars.AppName = p.Job
	jenkinsRenderVars.Custom = p.Custom
	rawConfigVars, err := mapz.DecodeStructToMap(jenkinsRenderVars)
	if err != nil {
		log.Debugf("jenkins config Jenkinsfile variables failed => %+v", err)
		return nil, err
	}
	ciConfig.Vars = rawConfigVars
	return ciConfig, nil
}
