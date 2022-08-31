package jenkins

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	jenkinsGitlabCredentialName = "jenkinsGitlabCredential"
	jenkinsGitlabConnectionName = "jenkinsGitlabConnection"
)

type JobOptions struct {
	JenkinsURL      string `mapstructure:"jenkinsURL" validate:"required,url"`
	JenkinsUser     string `mapstructure:"jenkinsUser"`
	JobName         string `mapstructure:"jobName"`
	JobFolderName   string `mapstructure:"jobFolderName"`
	ProjectURL      string `mapstructure:"projectURL" validate:"required"`
	ProjectBranch   string `mapstructure:"projectBranch"`
	JenkinsfilePath string `mapstructure:"jenkinsfilePath" validate:"required"`

	JenkinsEnableRestart bool `mapstructure:"jenkinsEnableRestart"`

	// used in package
	BasicAuth   *jenkins.BasicAuth `mapstructure:"basicAuth"`
	ProjectRepo *common.Repo       `mapstructure:"projectRepo"`
	CIConfig    *ci.CIConfig       `mapstructure:"ci"`
	SecretToken string             `mapstructure:"secretToken"`
}

type jobScriptRenderInfo struct {
	RepoType         string
	JobName          string
	RepositoryURL    string
	Branch           string
	SecretToken      string
	FolderName       string
	GitlabConnection string
}

func newJobOptions(options plugininstaller.RawOptions) (*JobOptions, error) {
	var opts JobOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (j *JobOptions) encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(j, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func (j *JobOptions) newJenkinsClient() (jenkins.JenkinsAPI, error) {
	return jenkins.NewClient(j.JenkinsURL, j.BasicAuth)
}

func (j *JobOptions) createOrUpdateJob(jenkinsClient jenkins.JenkinsAPI) error {
	// 1. render groovy script
	jobScript, err := jenkins.BuildRenderedScript(&jobScriptRenderInfo{
		RepoType:         j.ProjectRepo.RepoType,
		JobName:          j.JobName,
		RepositoryURL:    j.ProjectRepo.BuildURL(),
		Branch:           j.ProjectRepo.Branch,
		SecretToken:      j.SecretToken,
		FolderName:       j.JobFolderName,
		GitlabConnection: jenkinsGitlabConnectionName,
	})
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

func (j *JobOptions) buildWebhookInfo() *git.WebhookConfig {
	webHookURL := fmt.Sprintf("%s/project/%s", j.JenkinsURL, j.getJobPath())
	log.Debugf("jenkins config webhook is %s", webHookURL)
	return &git.WebhookConfig{
		Address:     webHookURL,
		SecretToken: j.SecretToken,
	}
}

func (j *JobOptions) installPlugins(jenkinsClient jenkins.JenkinsAPI, plugins []string) error {
	return jenkinsClient.InstallPluginsIfNotExists(plugins, j.JenkinsEnableRestart)
}

func (j *JobOptions) createGitlabConnection(jenkinsClient jenkins.JenkinsAPI, cascTemplate string) error {
	err := jenkinsClient.CreateGiltabCredential(jenkinsGitlabCredentialName, os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		log.Debugf("jenkins preinstall credentials failed: %s", err)
		return err
	}
	// 3. config gitlab casc
	cascConfig, err := template.Render(
		"jenkins-casc", cascTemplate, map[string]string{
			"CredentialName":       jenkinsGitlabCredentialName,
			"GitLabConnectionName": jenkinsGitlabConnectionName,
			"GitlabURL":            j.ProjectRepo.BaseURL,
		},
	)
	if err != nil {
		log.Debugf("jenkins preinstall credentials failed: %s", err)
		return err
	}
	return jenkinsClient.ConfigCasc(cascConfig)
}

func (j *JobOptions) deleteJob(client jenkins.JenkinsAPI) error {
	jobPath := j.getJobPath()
	if _, err := client.GetJob(context.Background(), jobPath); err == nil {
		if _, err := client.DeleteJob(context.Background(), jobPath); err != nil {
			return err
		}
	}
	return nil
}

func (j *JobOptions) getJobPath() string {
	return path.Join(j.JobFolderName, j.JobName)
}
