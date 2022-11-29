package configmanager

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	repoScaffoldingPluginName = "repo-scaffolding"
)

type repoTemplate struct {
	*scm.SCMInfo `yaml:",inline"`
	Vars         RawOptions `yaml:"vars"`
}

type app struct {
	Name         string        `yaml:"name" mapstructure:"name"`
	Spec         *appSpec      `yaml:"spec" mapstructure:"spec"`
	Repo         *scm.SCMInfo  `yaml:"repo" mapstructure:"repo"`
	RepoTemplate *repoTemplate `yaml:"repoTemplate" mapstructure:"repoTemplate"`
	CIRawConfigs []pipelineRaw `yaml:"ci" mapstructure:"ci"`
	CDRawConfigs []pipelineRaw `yaml:"cd" mapstructure:"cd"`
}

func (a *app) getTools(vars map[string]any, templateMap map[string]string) (Tools, error) {
	// generate app repo and template repo from scmInfo
	a.setDefault()
	repoScaffoldingTool, err := a.getRepoTemplateTool()
	if err != nil {
		return nil, fmt.Errorf("app[%s] can't get valid repo config: %w", a.Name, err)
	}

	// get ci/cd pipelineTemplates
	appVars := a.Spec.merge(vars)
	tools, err := a.generateCICDTools(templateMap, appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get pipeline tools failed: %w", a.Name, err)
	}
	if repoScaffoldingTool != nil {
		tools = append(tools, repoScaffoldingTool)
	}

	log.Debugf("Have got %d tools from app %s.", len(tools), a.Name)

	return tools, nil
}

// getAppPipelineTool generate ci/cd tools from app config
func (a *app) generateCICDTools(templateMap map[string]string, appVars map[string]any) (Tools, error) {
	allPipelineRaw := append(a.CIRawConfigs, a.CDRawConfigs...)
	var tools Tools
	for _, p := range allPipelineRaw {
		t, err := p.getPipelineTemplate(templateMap, appVars)
		if err != nil {
			return nil, err
		}
		pipelineTool, err := t.generatePipelineTool(a)
		if err != nil {
			return nil, err
		}
		pipelineTool.DependsOn = a.getRepoTemplateDependants()
		tools = append(tools, pipelineTool)
	}
	return tools, nil
}

// getRepoTemplateTool will use repo-scaffolding plugin for app
func (a *app) getRepoTemplateTool() (*Tool, error) {
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
		if a.RepoTemplate.Vars == nil {
			a.RepoTemplate.Vars = make(RawOptions)
		}
		return newTool(
			repoScaffoldingPluginName, a.Name, RawOptions{
				"destinationRepo": RawOptions(appRepo.Encode()),
				"sourceRepo":      RawOptions(templateRepo.Encode()),
				"vars":            a.RepoTemplate.Vars,
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
