package configmanager

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	repoScaffoldingPluginName = "repo-scaffolding"
)

type app struct {
	Name         string         `yaml:"name" mapstructure:"name"`
	Spec         map[string]any `yaml:"spec" mapstructure:"spec"`
	Repo         *scm.SCMInfo   `yaml:"repo" mapstructure:"repo"`
	RepoTemplate *scm.SCMInfo   `yaml:"repoTemplate" mapstructure:"repoTemplate"`
	CIRawConfigs []pipelineRaw  `yaml:"ci" mapstructure:"ci"`
	CDRawConfigs []pipelineRaw  `yaml:"cd" mapstructure:"cd"`
}

// getAppPipelineTool generate ci/cd tools from app config
func (a *app) generateCICDToolsFromAppConfig(templateMap map[string]string, appVars map[string]any) (Tools, error) {
	allPipelineRaw := append(a.CIRawConfigs, a.CDRawConfigs...)
	var tools Tools
	for _, p := range allPipelineRaw {
		t, err := p.newPipeline(a.Repo, templateMap, appVars)
		if err != nil {
			return nil, err
		}
		pipelineTool, err := t.getPipelineTool(a.Name)
		if err != nil {
			return nil, err
		}
		pipelineTool.DependsOn = a.getRepoTemplateDependants()
		tools = append(tools, *pipelineTool)
	}
	return tools, nil
}

// getRepoTemplateTool will use repo-scaffolding plugin for app
func (a *app) getRepoTemplateTool(appVars map[string]any) (*Tool, error) {
	if a.Repo == nil {
		return nil, fmt.Errorf("app.repo field can't be empty")
	}
	appRepo, err := a.Repo.BuildRepoInfo()
	if err != nil {
		return nil, fmt.Errorf("configmanager[app] parse repo failed: %w", err)
	}
	if a.RepoTemplate != nil {
		templateRepo, err := a.RepoTemplate.BuildRepoInfo()
		if err != nil {
			return nil, fmt.Errorf("configmanager[app] parse repoTemplate failed: %w", err)
		}
		return newTool(
			repoScaffoldingPluginName, a.Name, RawOptions{
				"destinationRepo": RawOptions(appRepo.Encode()),
				"sourceRepo":      RawOptions(templateRepo.Encode()),
				"vars":            RawOptions(appVars),
			},
		), nil
	}
	return nil, nil
}

// setDefault will set repoName to appName if repo.name field is empty
func (a *app) setDefault() {
	if a.Repo != nil && a.Repo.Name == "" {
		a.Repo.Name = a.Name
	}
}

// since all plugin depends on code is deployed, get dependsOn for repoTemplate
func (a *app) getRepoTemplateDependants() []string {
	var dependsOn []string
	// if a.RepoTemplate is configured, pipeline need to wait reposcaffolding finished
	if a.RepoTemplate != nil {
		dependsOn = []string{fmt.Sprintf("%s.%s", repoScaffoldingPluginName, a.Name)}
	}
	return dependsOn
}
