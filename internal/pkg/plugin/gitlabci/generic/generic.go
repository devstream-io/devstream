package generic

import (
	_ "embed"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
	"github.com/devstream-io/devstream/pkg/util/template"
)

//go:embed tpl/helmValue.tpl.yaml
var helmValueTpl string

var ciType = server.CIGitLabType

func preConfigGitlab(options configmanager.RawOptions) error {
	opts, err := ci.NewCIOptions(options)
	if err != nil {
		return err
	}

	// PreConfig steps
	stepConfigs := step.ExtractValidStepConfig(opts.Pipeline)
	scmClient, err := scm.NewClientWithAuth(opts.ProjectRepo)
	if err != nil {
		log.Debugf("init gitlab client failed: %+v", err)
		return err
	}
	for _, c := range stepConfigs {
		err := c.ConfigSCM(scmClient)
		if err != nil {
			log.Debugf("gitlabci config step failed: %+v", err)
			return err
		}
	}

	// push shareLib to gitlab
	if err := repoScaffoldShareLib(opts.ProjectRepo); err != nil {
		return err
	}
	return createGitlabRunnerByHelm(opts.ProjectRepo)
}

func repoScaffoldShareLib(repoInfo *git.RepoInfo) error {
	repoScaffoldOptions := configmanager.RawOptions{
		"sourceRepo": map[string]interface{}{
			"owner":   "devstream-io",
			"name":    "gitlab-ci-library",
			"scmType": "github",
		},
		"destinationRepo": map[string]interface{}{
			"owner":   repoInfo.GetRepoOwner(),
			"branch":  "master",
			"name":    "gitlab-ci-library",
			"scmType": "gitlab",
			"baseURL": repoInfo.BaseURL,
		},
	}
	return reposcaffolding.InstallRepo(repoScaffoldOptions)
}

func createGitlabRunnerByHelm(repoInfo *git.RepoInfo) error {
	gitlabClient, err := gitlab.NewClient(repoInfo)
	if err != nil {
		log.Debugf("gitlabci init gitlab client failed: %+v", err)
		return err
	}
	// 1. check project has runner, if project has runner, just return
	runner, err := gitlabClient.ListRepoRunner()
	if err != nil {
		log.Debugf("gitlabci get project runner failed: %+v", err)
		return err
	}
	if len(runner) > 0 {
		log.Debugf("gitlabci runner exist")
		return nil
	}
	runnerToken, err := gitlabClient.ResetRepoRunnerToken()
	if err != nil {
		log.Debugf("gitlabci reset project runner token failed: %+v", err)
		return err
	}
	// 2. else create runner for this project
	valuesYaml, err := template.New().FromContent(helmValueTpl).SetDefaultRender(
		"gitlab-runner tpl", map[string]string{
			"GitlabURL":     repoInfo.BaseURL,
			"RegisterToken": runnerToken,
		},
	).Render()
	if err != nil {
		return err
	}
	helmOptions := configmanager.RawOptions{
		"instanceID": "gitlab-global-runner",
		"repo": map[string]interface{}{
			"name": "gitlab",
			"url":  "https://charts.gitlab.io",
		},
		"chart": map[string]interface{}{
			"chartPath":   "",
			"chartName":   "gitlab/gitlab-runner",
			"ValuesYaml":  valuesYaml,
			"wait":        true,
			"timeout":     "10m",
			"upgradeCRDs": false,
			"namespace":   "gitlab",
			"releaseName": fmt.Sprintf("%s-runner", repoInfo.Repo),
		},
	}
	_, err = helminstaller.Create(helmOptions)
	return err
}
