package jenkins

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type JobOptions struct {
	Jenkins  Jenkins  `mapstructure:"jenkins"`
	SCM      SCM      `mapstructure:"scm"`
	Pipeline Pipeline `mapstructure:"pipeline"`

	// used in package
	BasicAuth   *jenkins.BasicAuth `mapstructure:"basicAuth"`
	ProjectRepo *common.Repo       `mapstructure:"projectRepo"`
	CIConfig    *ci.CIConfig       `mapstructure:"ci"`
	SecretToken string             `mapstructure:"secretToken"`
}

type Jenkins struct {
	URL           string `mapstructure:"url" validate:"required,url"`
	User          string `mapstructure:"user"`
	Namespace     string `mapstructure:"namespace"`
	EnableRestart bool   `mapstructure:"enableRestart"`
}

type SCM struct {
	CloneURL string `mapstructure:"cloneURL" validate:"required"`
	APIURL   string `mapstructure:"apiURL"`
	Branch   string `mapstructure:"branch"`
	Type     string `mapstructure:"type"`

	// used in package
	SSHprivateKey string `mapstructure:"sshPrivateKey"`
}

type jobScriptRenderInfo struct {
	RepoType          string
	JobName           string
	RepositoryURL     string
	Branch            string
	SecretToken       string
	FolderName        string
	GitlabConnection  string
	RepoCredentialsId string
	RepoURL           string
	RepoName          string
	RepoOwner         string
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
	return jenkins.NewClient(j.Jenkins.URL, j.BasicAuth)
}

func (j *JobOptions) buildWebhookInfo() *git.WebhookConfig {
	var webHookURL string
	switch j.ProjectRepo.RepoType {
	case "gitlab":
		webHookURL = fmt.Sprintf("%s/project/%s", j.Jenkins.URL, j.Pipeline.getJobPath())
	case "github":
		webHookURL = fmt.Sprintf("%s/github-webhook/", j.Jenkins.URL)
	}
	log.Debugf("jenkins config webhook is %s", webHookURL)
	return &git.WebhookConfig{
		Address:     webHookURL,
		SecretToken: j.SecretToken,
	}
}

func (j *JobOptions) deleteJob(client jenkins.JenkinsAPI) error {
	job, err := client.GetFolderJob(j.Pipeline.getJobName(), j.Pipeline.getJobFolder())
	if err != nil {
		// check job is already been deleted
		if jenkins.IsNotFoundError(err) {
			return err
		}
		return nil
	}
	isDeleted, err := job.Delete(context.Background())
	if err != nil {
		log.Debugf("jenkins delete job %s failed: %s", j.Pipeline.getJobPath(), err)
		return err
	}
	log.Debugf("jenkins delete job %s status: %v", j.Pipeline.getJobPath(), isDeleted)
	return nil
}

func (j *JobOptions) buildCIConfig() {
	jenkinsFilePath := j.Pipeline.JenkinsfilePath
	ciConfig := &ci.CIConfig{
		Type: "jenkins",
	}
	// config CIConfig
	jenkinsfileURL, err := url.ParseRequestURI(jenkinsFilePath)
	// if path is url, download from remote
	if err != nil || jenkinsfileURL.Host == "" {
		ciConfig.LocalPath = jenkinsFilePath
	} else {
		ciConfig.RemoteURL = jenkinsFilePath
	}
	var imageName string
	if j.ProjectRepo != nil {
		imageName = j.ProjectRepo.Repo
	} else {
		imageName = j.Pipeline.JobName
	}
	harborURLHost := path.Join(j.Pipeline.getImageHost(), defaultImageProject)
	ciConfig.Vars = map[string]interface{}{
		"ImageName":        imageName,
		"ImageRepoAddress": harborURLHost,
	}
	j.CIConfig = ciConfig
}

func (j *JobOptions) extractJenkinsPlugins() []pluginConfigAPI {
	var pluginsConfigs []pluginConfigAPI
	switch j.ProjectRepo.RepoType {
	case "gitlab":
		pluginsConfigs = append(pluginsConfigs, &plugins.GitlabJenkinsConfig{
			SSHPrivateKey: j.SCM.SSHprivateKey,
			RepoOwner:     j.ProjectRepo.Owner,
			BaseURL:       j.ProjectRepo.BaseURL,
		})
	case "github":
		pluginsConfigs = append(pluginsConfigs, &plugins.GithubJenkinsConfig{
			JenkinsURL: j.Jenkins.URL,
		})
	}
	return pluginsConfigs
}

func (j *JobOptions) createOrUpdateJob(jenkinsClient jenkins.JenkinsAPI) error {
	// 1. render groovy script
	jobRenderInfo := &jobScriptRenderInfo{
		RepoType:         j.ProjectRepo.RepoType,
		JobName:          j.Pipeline.getJobName(),
		RepositoryURL:    j.SCM.CloneURL,
		Branch:           j.ProjectRepo.Branch,
		SecretToken:      j.SecretToken,
		FolderName:       j.Pipeline.getJobFolder(),
		GitlabConnection: plugins.GitlabConnectionName,
		RepoURL:          j.ProjectRepo.BuildURL(),
		RepoOwner:        j.ProjectRepo.Owner,
		RepoName:         j.ProjectRepo.Repo,
	}
	// config credential for different repo type
	switch j.ProjectRepo.RepoType {
	case "gitlab":
		if j.SCM.SSHprivateKey != "" {
			jobRenderInfo.RepoCredentialsId = plugins.SSHKeyCredentialName
		}
	case "github":
		jobRenderInfo.RepoCredentialsId = plugins.GithubCredentialName
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
