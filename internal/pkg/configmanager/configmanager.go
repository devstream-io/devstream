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
// 2. Validation.
func (m *Manager) LoadConfig() (*Config, error) {
	// step 1: get config
	c, err := m.getConfigFromFileWithGlobalVars()
	if err != nil {
		return nil, err
	}
	// set instanceID in options
	c.renderInstanceIDtoOptions()

	// step 2: check config is valid
	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// getConfigFromFileWithGlobalVars gets Config from the config file specified by Manager.ConfigFilePath, then:
// 1. render the global variables to Config.Tools and Config.Apps
// 2. transfer the PipelineTemplates to Config.pipelineTemplateMap, it's map[string]string type.
// We can't render the original config file to Config.PipelineTemplates directly for the:
//  1. variables rendered must be before the yaml.Unmarshal() called for the [[ foo ]] will be treated as a two-dimensional array by the yaml parser;
//  2. the variables used([[ foo ]]) in the Config.PipelineTemplates can be defined in the Config.Apps or Config.Vars;
func (m *Manager) getConfigFromFileWithGlobalVars() (*Config, error) {
	configBytes, err := os.ReadFile(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// extract top raw config struct from config text
	r, err := newRawConfigFromConfigBytes(configBytes)
	if err != nil {
		return nil, err
	}
	// 1. get global variables
	vars, err := r.getVars()
	if err != nil {
		return nil, fmt.Errorf("failed to get variables from config file. Error: %w", err)
	}

	// 2. tools with global variables rendered
	tools, err := r.getToolsWithVars(vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools from config file. Error: %w", err)
	}

	// 3. apps tools with global variables rendered
	appTools, err := r.getAppToolsWithVars(vars)
	if err != nil {
		return nil, fmt.Errorf("failed to get apps from config file. Error: %w", err)
	}
	// all tools from apps should depend on the original tools,
	// because dtm will execute all original tools first, then execute all tools from apps
	appTools.updateToolDepends(tools)
	tools = append(tools, appTools...)

	// 4. coreConfig without any changes
	coreConfig, err := r.getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get coreConfig from config file. Error: %w", err)
	}

	return &Config{
		Config: *coreConfig,
		Vars:   vars,
		Tools:  tools,
	}, nil
}
