package configmanager

import (
	"fmt"
	"os"
)

// Manager is used to load the config file from the ConfigFilePath and finally get the Config object.
type Manager struct {
	ConfigFilePath string
}

// NewManager takes configFilePath, then return a *Manager object.
func NewManager(configFilePath string) *Manager {
	return &Manager{
		ConfigFilePath: configFilePath,
	}
}

// LoadConfig is the only method that the caller of Manager needs to be concerned with, and this method returns a *Config finally.
// The main workflow of this method is:
// 1. Get the original config from the config file specified by ConfigFilePath;
// 2. Parsing tools from the apps, and merge that tools to Config.Tools;
// 3. Validation.
func (m *Manager) LoadConfig() (*Config, error) {
	// step 1
	c, err := m.getConfigFromFileWithGlobalVars()
	if err != nil {
		return nil, err
	}

	// step 2
	appTools, err := c.getToolsFromApps()
	if err != nil {
		return nil, err
	}

	c.Tools = append(c.Tools, appTools...)
	c.renderInstanceIDtoOptions()

	// step 3
	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// getConfigFromFileWithGlobalVars gets Config from the config file specified by Manager.ConfigFilePath, then:
// 1. render the global variables to Config.Tools and Config.Apps
// 2. transfer the PipelineTemplates to Config.pipelineTemplateMap, it's map[string]string type.
// We can't render the original config file to Config.PipelineTemplates directly for the:
//   1. variables rendered must be before the yaml.Unmarshal() called for the [[ foo ]] will be treated as a two-dimensional array by the yaml parser;
//   2. the variables used([[ foo ]]) in the Config.PipelineTemplates can be defined in the Config.Apps or Config.Vars;
//   3. pipeline templates are used in apps, so it would be more appropriate to refer to pipeline templates when dealing with apps
func (m *Manager) getConfigFromFileWithGlobalVars() (*Config, error) {
	configBytes, err := os.ReadFile(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// global variables
	vars, err := getVarsFromConfigFile(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get variables from config file. Error: %w", err)
	}

	// tools with global variables rendered
	tools, err := getToolsFromConfigFileWithVarsRendered(configBytes, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools from config file. Error: %w", err)
	}

	// apps with global variables rendered
	apps, err := getAppsFromConfigFileWithVarsRendered(configBytes, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get apps from config file. Error: %w", err)
	}

	// pipelineTemplateMap transfer from PipelineTemplates
	pipelineTemplateMap, err := getPipelineTemplatesMapFromConfigFile(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipelineTemplatesMap from config file. Error: %w", err)
	}

	// coreConfig without any changes
	coreConfig, err := getCoreConfigFromConfigFile(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get coreConfig from config file. Error: %w", err)
	}

	return &Config{
		Config:              *coreConfig,
		Vars:                vars,
		Tools:               tools,
		Apps:                apps,
		pipelineTemplateMap: pipelineTemplateMap,
	}, nil
}
