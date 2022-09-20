package jenkins

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	jenkinsGitlabCredentialName = "jenkinsGitlabCredential"
	jenkinsGitlabConnectionName = "jenkinsGitlabConnection"
	jenkinsSSHkeyCredentialName = "jenkinsSSHKeyCredential"
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
	CloneURL      string `mapstructure:"cloneURL" validate:"required"`
	APIURL        string `mapstructure:"apiURL"`
	Branch        string `mapstructure:"branch"`
	Type          string `mapstructure:"type"`
	SSHprivateKey string `mapstructure:"sshPrivateKey"`
}

type Pipeline struct {
	JobName         string    `mapstructure:"jobName" validate:"required"`
	JenkinsfilePath string    `mapstructure:"jenkinsfilePath" validate:"required"`
	ImageRepo       ImageRepo `mapstructure:"imageRepo"`
}

type ImageRepo struct {
	URL  string `mapstructure:"url" validate:"url"`
	User string `mapstructure:"user"`
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

func (j *JobOptions) createOrUpdateJob(jenkinsClient jenkins.JenkinsAPI) error {
	// 1. render groovy script
	jobRenderInfo := &jobScriptRenderInfo{
		RepoType:         j.ProjectRepo.RepoType,
		JobName:          j.getJobName(),
		RepositoryURL:    j.SCM.CloneURL,
		Branch:           j.ProjectRepo.Branch,
		SecretToken:      j.SecretToken,
		FolderName:       j.getJobFolder(),
		GitlabConnection: jenkinsGitlabConnectionName,
	}
	if j.SCM.SSHprivateKey != "" {
		jobRenderInfo.RepoCredentialsId = jenkinsSSHkeyCredentialName
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

func (j *JobOptions) buildWebhookInfo() *git.WebhookConfig {
	webHookURL := fmt.Sprintf("%s/project/%s", j.Jenkins.URL, j.getJobPath())
	log.Debugf("jenkins config webhook is %s", webHookURL)
	return &git.WebhookConfig{
		Address:     webHookURL,
		SecretToken: j.SecretToken,
	}
}

func (j *JobOptions) installPlugins(jenkinsClient jenkins.JenkinsAPI, plugins []string) error {
	return jenkinsClient.InstallPluginsIfNotExists(plugins, j.Jenkins.EnableRestart)
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
	return j.Pipeline.JobName
}

func (j *JobOptions) getJobFolder() string {
	if strings.Contains(j.Pipeline.JobName, "/") {
		return strings.Split(j.Pipeline.JobName, "/")[0]
	}
	return ""
}

func (j *JobOptions) getJobName() string {
	if strings.Contains(j.Pipeline.JobName, "/") {
		return strings.Split(j.Pipeline.JobName, "/")[1]
	}
	return j.Pipeline.JobName
}

func (j *JobOptions) getImageHost() string {
	harborAddress := j.Pipeline.ImageRepo.URL
	harborURL, err := url.ParseRequestURI(harborAddress)
	if err != nil {
		return harborAddress
	}
	return harborURL.Host
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
	harborURLHost := path.Join(j.getImageHost(), defaultImageProject)
	ciConfig.Vars = map[string]interface{}{
		"ImageName":        imageName,
		"ImageRepoAddress": harborURLHost,
	}
	j.CIConfig = ciConfig
}

func (j *JobOptions) createGitlabSSHPrivateKey(jenkinsClient jenkins.JenkinsAPI) error {
	if j.SCM.SSHprivateKey == "" {
		log.Warnf("jenkins gitlab ssh key not config, private repo can't be clone")
		return nil
	}
	return jenkinsClient.CreateSSHKeyCredential(
		jenkinsSSHkeyCredentialName, j.ProjectRepo.Owner, j.SCM.SSHprivateKey,
	)
}
