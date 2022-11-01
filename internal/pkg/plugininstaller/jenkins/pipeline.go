package jenkins

import (
	"context"
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type pipeline struct {
	Job             string                    `mapstructure:"jobName" validate:"required"`
	JenkinsfilePath string                    `mapstructure:"jenkinsfilePath" validate:"required"`
	ImageRepo       *step.ImageRepoStepConfig `mapstructure:"imageRepo"`
	Dingtalk        *step.DingtalkStepConfig  `mapstructure:"dingTalk"`
	Sonarqube       *step.SonarQubeStepConfig `mapstructure:"sonarqube"`
	General         *step.GeneralStepConfig   `mapstructure:"general"`
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

func (p *pipeline) extractPlugins(repoInfo *git.RepoInfo) []step.StepConfigAPI {
	stepConfigs := step.ExtractValidStepConfig(p)
	// add repo plugin for repoInfo
	stepConfigs = append(stepConfigs, step.GetRepoStepConfig(repoInfo)...)
	return stepConfigs
}

func (p *pipeline) install(jenkinsClient jenkins.JenkinsAPI, repoInfo *git.RepoInfo, secretToken string) error {
	// 1. install jenkins plugins
	pipelinePlugins := p.extractPlugins(repoInfo)
	if err := ensurePluginInstalled(jenkinsClient, pipelinePlugins); err != nil {
		return err
	}
	// 2. config plugins by casc
	if err := configPlugins(jenkinsClient, pipelinePlugins); err != nil {
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
	globalConfig := step.GetStepGlobalVars(repoInfo)
	// 1. render groovy script
	jobRenderInfo := &jenkins.JobScriptRenderInfo{
		RepoType:          repoInfo.RepoType,
		JobName:           p.getJobName(),
		RepositoryURL:     repoInfo.CloneURL,
		Branch:            repoInfo.Branch,
		SecretToken:       secretToken,
		FolderName:        p.getJobFolder(),
		GitlabConnection:  globalConfig.GitlabConnectionID,
		RepoURL:           repoInfo.BuildScmURL(),
		RepoOwner:         repoInfo.GetRepoOwner(),
		RepoName:          repoInfo.Repo,
		RepoCredentialsId: globalConfig.CredentialID,
	}
	jobScript, err := jenkins.BuildRenderedScript(jobRenderInfo)
	if err != nil {
		log.Debugf("jenkins redner template failed: %s", err)
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

func (p *pipeline) buildCIConfig(repoInfo *git.RepoInfo, pipelineRawOption map[string]interface{}) *cifile.CIConfig {
	ciConfig := &cifile.CIConfig{
		Type:           ciType,
		ConfigLocation: p.JenkinsfilePath,
	}
	// update ci render variables by plugins
	rawConfigVars := step.GenerateCIFileVars(p, repoInfo)
	rawConfigVars.Set("AppName", p.Job)
	ciConfig.Vars = rawConfigVars
	log.Debugf("jenkins pipeline get render vars: %+v", ciConfig.Vars)
	return ciConfig
}
