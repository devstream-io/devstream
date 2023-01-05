package gitlabci

import (
	_ "embed"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/helminstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
	"github.com/devstream-io/devstream/pkg/util/template"
)

//go:embed tpl/helmValue.tpl.yaml
var helmValueTpl string

func preConfigGitlab(rawOptions configmanager.RawOptions) error {
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
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
	// create runner in self-host gitlab
	if opts.needCreateRunner() {
		return createGitlabRunnerByHelm(opts.ProjectRepo)
	}
	return nil
}

// createGitlabRunnerByHelm will install gitlab runner if it's not exist
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
	valuesYaml, err := template.NewRenderClient(&template.TemplateOption{
		Name: "gitlab-runner"}, template.ContentGetter).Render(
		helmValueTpl, map[string]any{
			"GitlabURL":     repoInfo.BaseURL,
			"RegisterToken": runnerToken,
		})
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
	if err != nil {
		log.Debugf("gitlab ci install runner by helm failed: %s", err)
		return fmt.Errorf("gitlabci create runner by helm failed, please check your kubeneretes config")
	}
	// disable gitlab shared runner
	return gitlabClient.DisableRepoSharedRunner()
}
