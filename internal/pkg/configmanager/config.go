package configmanager

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// Config is a general config in DevStream.
type Config struct {
	Config CoreConfig     `yaml:"config"`
	Vars   map[string]any `yaml:"vars"`
	Tools  Tools          `yaml:"tools"`
	Apps   []*app         `yaml:"apps"`
	// We'll read the pipeline templates from config file and render it to pipelineTemplateMap in no time.
	//PipelineTemplates   []*pipelineTemplate `yaml:"-"`
	pipelineTemplateMap map[string]string `yaml:"-"`
}

func (c *Config) getToolsFromApps() (Tools, error) {
	var tools Tools
	for _, a := range c.Apps {
		appTools, err := c.getToolsFromApp(a)
		if err != nil {
			return nil, err
		}
		tools = append(tools, appTools...)
	}
	return tools, nil
}

func (c *Config) getToolsFromApp(a *app) (Tools, error) {
	// generate app repo and template repo from scmInfo
	a.setDefault()
	repoScaffoldingTool, err := a.getRepoTemplateTool()
	if err != nil {
		return nil, fmt.Errorf("app[%s] get repo failed: %w", a.Name, err)
	}

	// get ci/cd pipelineTemplates
	appVars := a.Spec.merge(c.Vars)
	tools, err := a.generateCICDToolsFromAppConfig(c.pipelineTemplateMap, appVars)
	if err != nil {
		return nil, fmt.Errorf("app[%s] get pipeline tools failed: %w", a.Name, err)
	}
	if repoScaffoldingTool != nil {
		tools = append(tools, repoScaffoldingTool)
	}

	// all tools from apps should depend on the original tools,
	// because dtm will execute all original tools first, then execute all tools from apps
	for _, toolFromApps := range tools {
		for _, t := range c.Tools {
			toolFromApps.DependsOn = append(toolFromApps.DependsOn, t.KeyWithNameAndInstanceID())
		}
	}

	log.Debugf("Have got %d tools from app %s.", len(tools), a.Name)
	for i, t := range tools {
		log.Debugf("Tool %d: %v", i+1, t)
	}

	return tools, nil
}

func (c *Config) renderInstanceIDtoOptions() {
	for _, t := range c.Tools {
		if t.Options == nil {
			t.Options = make(RawOptions)
		}
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
