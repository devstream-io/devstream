package jenkins

import (
	"fmt"
	"path"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
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
}

type jobScriptRenderInfo struct {
	RepoType      string
	JobName       string
	RepositoryURL string
	Branch        string
	SecretToken   string
	FolderName    string
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

func (j *JobOptions) buildRenderOptions(secretToekn string) *jobScriptRenderInfo {
	return &jobScriptRenderInfo{
		RepoType:      j.ProjectRepo.RepoType,
		JobName:       j.JobName,
		RepositoryURL: j.ProjectRepo.BuildURL(),
		Branch:        j.ProjectRepo.Branch,
		SecretToken:   secretToekn,
		FolderName:    j.JobFolderName,
	}
}

func (j *JobOptions) buildWebhookInfo(secretToken string) *git.WebhookConfig {
	webHookURL := fmt.Sprintf("%s/project/%s", j.JenkinsURL, path.Join(j.JobFolderName, j.JobName))
	log.Debugf("jenkins config webhook is %s", webHookURL)
	return &git.WebhookConfig{
		Address:     webHookURL,
		SecretToken: secretToken,
	}
}
