package configmanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
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
	err := c.renderPipelineTemplateMap()
	if err != nil {
		return nil, err
	}

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

	// 1. render appStr with Vars
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

	// 3. generate app repo and template repo from scmInfo
	rawApp.setDefault()
	repoScaffoldingTool, err := rawApp.getRepoTemplateTool()
	if err != nil {
		return nil, fmt.Errorf("app[%s] get repo failed: %w", rawApp.Name, err)
	}

	// 4. get ci/cd pipelineTemplates
	appVars := rawApp.Spec.merge(c.Vars)
	tools, err := rawApp.generateCICDToolsFromAppConfig(c.pipelineTemplateMap, appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get pipeline tools failed: %w", rawApp.Name, err)
	}
	if repoScaffoldingTool != nil {
		tools = append(tools, repoScaffoldingTool)
	}

	log.Debugf("Have got %d tools from app %s.", len(tools), rawApp.Name)
	for i, t := range tools {
		log.Debugf("Tool %d: %v", i+1, t)
	}

	return tools, nil
}

func (c *Config) renderToolsWithVars() error {
	toolsStr, err := yaml.Marshal(c.Tools)
	if err != nil {
		return err
	}

	toolsStrWithVars, err := renderConfigWithVariables(string(toolsStr), c.Vars)
	if err != nil {
		return err
	}

	var tools Tools
	if err = yaml.Unmarshal(toolsStrWithVars, &tools); err != nil {
		return err
	}
	c.Tools = tools

	return nil
}

func (c *Config) renderPipelineTemplateMap() error {
	c.pipelineTemplateMap = make(map[string]string)
	for _, pt := range c.PipelineTemplates {
		tplBytes, err := yaml.Marshal(pt)
		if err != nil {
			return err
		}
		c.pipelineTemplateMap[pt.Name] = string(tplBytes)
	}
	return nil
}

func (c *Config) renderInstanceIDtoOptions() {
	for _, t := range c.Tools {
		t.Options["instanceID"] = t.InstanceID
	}
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

func (c *Config) String() string {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}
