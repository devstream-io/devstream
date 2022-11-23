package configmanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
)

// Config is a general config in DevStream.
type Config struct {
	Config              CoreConfig         `yaml:"config"`
	Vars                map[string]any     `yaml:"vars"`
	Tools               Tools              `yaml:"tools"`
	Apps                []app              `yaml:"apps"`
	PipelineTemplates   []pipelineTemplate `yaml:"pipelineTemplates"`
	pipelineTemplateMap map[string]string  `yaml:"-"`
}

func (c *Config) getToolsFromApps() (Tools, error) {
	var tools Tools
	for _, a := range c.Apps {
		appTools, err := c.getToolsWithVarsFromApp(a)
		if err != nil {
			return nil, err
		}
		tools = append(tools, appTools...)
	}
	return tools, nil
}

func (c *Config) getToolsWithVarsFromApp(a app) (Tools, error) {
	appStr, err := yaml.Marshal(a)
	if err != nil {
		return nil, err
	}

	// 1. render appStr with globalVars
	appRenderStr, err := renderConfigWithVariables(string(appStr), c.Vars)
	if err != nil {
		log.Debugf("configmanager/app %s render globalVars %+v failed", appRenderStr, c.Vars)
		return nil, fmt.Errorf("app render globalVars failed: %w", err)
	}

	// 2. unmarshal app config for render pipelineTemplate
	var rawApp app
	if err = yaml.Unmarshal(appRenderStr, &rawApp); err != nil {
		return nil, fmt.Errorf("app parse yaml failed: %w", err)
	}

	err = c.renderPipelineTemplateMap()
	if err != nil {
		return nil, err
	}

	rawApp.setDefault()
	appVars := mapz.Merge(c.Vars, rawApp.Spec)

	// 3. generate app repo and tempalte repo from scmInfo
	repoScaffoldingTool, err := rawApp.getRepoTemplateTool(appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get repo failed: %w", rawApp.Name, err)
	}

	// 4. get ci/cd pipelineTemplates
	tools, err := rawApp.generateCICDToolsFromAppConfig(c.pipelineTemplateMap, appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get pipeline tools failed: %w", rawApp.Name, err)
	}
	if repoScaffoldingTool != nil {
		tools = append(tools, *repoScaffoldingTool)
	}

	return tools, nil
}

func (c *Config) renderPipelineTemplateMap() error {
	templateMap := make(map[string]string)
	for _, pt := range c.PipelineTemplates {
		tplBytes, err := yaml.Marshal(pt)
		if err != nil {
			return err
		}
		templateMap[pt.Name] = string(tplBytes)
	}
	return nil
}

func (c *Config) validate() error {
	if c.Config.State == nil {
		return fmt.Errorf("state is not defined")
	}

	if err := c.Tools.validateAll(); err != nil {
		return err
	}
	return nil
}
