package configmanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	repoScaffoldingPluginName = "repo-scaffolding"
)

type RawApp struct {
	Name         string         `yaml:"name" mapstructure:"name"`
	Spec         map[string]any `yaml:"spec" mapstructure:"spec"`
	Repo         *scm.SCMInfo   `yaml:"repo" mapstructure:"repo"`
	RepoTemplate *scm.SCMInfo   `yaml:"repoTemplate" mapstructure:"repoTemplate"`
	CIRawConfigs []pipelineRaw  `yaml:"ci" mapstructure:"ci"`
	CDRawConfigs []pipelineRaw  `yaml:"cd" mapstructure:"cd"`
}

// getToolsFromApp return app tools
func getToolsFromApp(appStr string, globalVars map[string]any, templateMap map[string]string) (Tools, error) {
	//1. render appStr with globalVars
	appRenderStr, err := renderConfigWithVariables(appStr, globalVars)
	if err != nil {
		log.Debugf("configmanager/app %s render globalVars %+v failed", appRenderStr, globalVars)
		return nil, fmt.Errorf("app render globalVars failed: %w", err)
	}
	// 2. unmarshal RawApp config for render pipelineTemplate
	var rawData RawApp
	if err := yaml.Unmarshal([]byte(appRenderStr), &rawData); err != nil {
		return nil, fmt.Errorf("app parse yaml failed: %w", err)
	}
	rawData.setDefault()
	appVars := mapz.Merge(globalVars, rawData.Spec)
	// 3. generate app repo and tempalte repo from scmInfo
	repoScaffoldingTool, err := rawData.getRepoTemplateTool(appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get repo failed: %w", rawData.Name, err)
	}
	// 4. get ci/cd pipelineTemplates
	tools, err := rawData.generateCICDToolsFromAppConfig(templateMap, appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get pipeline tools failed: %w", rawData.Name, err)
	}
	if repoScaffoldingTool != nil {
		tools = append(tools, *repoScaffoldingTool)
	}
	return tools, nil
}

// getAppPipelineTool generate ci/cd tools from app config
func (a *RawApp) generateCICDToolsFromAppConfig(templateMap map[string]string, appVars map[string]any) (Tools, error) {
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
func (a *RawApp) getRepoTemplateTool(appVars map[string]any) (*Tool, error) {
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
func (a *RawApp) setDefault() {
	if a.Repo != nil && a.Repo.Name == "" {
		a.Repo.Name = a.Name
	}
}

// since all plugin depends on code is deployed, get dependsOn for repoTemplate
func (a *RawApp) getRepoTemplateDependants() []string {
	var dependsOn []string
	// if a.RepoTemplate is configured, pipeline need to wait reposcaffolding finished
	if a.RepoTemplate != nil {
		dependsOn = []string{fmt.Sprintf("%s.%s", repoScaffoldingPluginName, a.Name)}
	}
	return dependsOn
}
