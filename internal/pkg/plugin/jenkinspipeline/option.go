package jenkinspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	ciType = server.CIServerType("jenkins")
)

type jenkinsJobName string

type jobOptions struct {
	Jenkins     jenkinsOption `mapstructure:"jenkins"`
	ci.CIConfig `mapstructure:",squash"`
	JobName     jenkinsJobName `mapstructure:"jobName"`
	// SCM      scm.SCMInfo   `mapstructure:"scm"`
	// Pipeline pipeline      `mapstructure:"pipeline"`

	// used in package
	// CIFileConfig *cifile.CIFileConfig `mapstructure:"ci"`
	// ProjectRepo  *git.RepoInfo        `mapstructure:"projectRepo"`
}

func newJobOptions(options configmanager.RawOptions) (*jobOptions, error) {
	var opts jobOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (n jenkinsJobName) getJobName() string {
	jobNameStr := string(n)
	if strings.Contains(jobNameStr, "/") {
		return strings.Split(jobNameStr, "/")[1]
	}
	return jobNameStr
}

func (n jenkinsJobName) getJobFolder() string {
	jobNameStr := string(n)
	if strings.Contains(jobNameStr, "/") {
		return strings.Split(jobNameStr, "/")[0]
	}
	return ""
}

func (n jenkinsJobName) checkValid() error {
	jobNameStr := string(n)
	if strings.Contains(jobNameStr, "/") {
		strs := strings.Split(jobNameStr, "/")
		if len(strs) != 2 || len(strs[0]) == 0 || len(strs[1]) == 0 {
			return fmt.Errorf("jenkins jobName illegal: %s", n)
		}
	}
	return nil
}

func (j *jobOptions) extractPlugins() []step.StepConfigAPI {
	stepConfigs := step.ExtractValidStepConfig(j.Pipeline)
	// add repo plugin for repoInfo
	stepConfigs = append(stepConfigs, step.GetRepoStepConfig(j.ProjectRepo)...)
	return stepConfigs
}

func (j *jobOptions) remove(jenkinsClient jenkins.JenkinsAPI) error {
	job, err := jenkinsClient.GetFolderJob(j.JobName.getJobName(), j.JobName.getJobFolder())
	if err != nil {
		// check job is already been deleted
		if jenkins.IsNotFoundError(err) {
			return nil
		}
		return err
	}
	_, err = job.Delete(context.Background())
	log.Debugf("jenkins delete job %s status: %v", j.JobName, err)
	return nil
}

func (j *jobOptions) createOrUpdateJob(jenkinsClient jenkins.JenkinsAPI, secretToken string) error {
	repoInfo := j.ProjectRepo
	globalConfig := step.GetStepGlobalVars(repoInfo)
	// 1. render groovy script
	jobRenderInfo := &jenkins.JobScriptRenderInfo{
		RepoType:          repoInfo.RepoType,
		JobName:           j.JobName.getJobName(),
		RepositoryURL:     repoInfo.CloneURL,
		Branch:            repoInfo.Branch,
		SecretToken:       secretToken,
		FolderName:        j.JobName.getJobFolder(),
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

func (j *jobOptions) install(jenkinsClient jenkins.JenkinsAPI, secretToken string) error {
	// 1. install jenkins plugins
	pipelinePlugins := j.extractPlugins()
	if err := ensurePluginInstalled(jenkinsClient, pipelinePlugins); err != nil {
		return err
	}
	// 2. config plugins by casc
	if err := configPlugins(jenkinsClient, pipelinePlugins); err != nil {
		return err
	}
	// 3. create or update jenkins job
	return j.createOrUpdateJob(jenkinsClient, secretToken)
}
