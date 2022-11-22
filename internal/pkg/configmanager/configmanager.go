package configmanager

type Manager struct {
	ConfigFilePath string
}

func NewManager(configFilePath string) *Manager {
	return &Manager{
		ConfigFilePath: configFilePath,
	}
}

// LoadConfig reads an input file as a general config.
// It will return "non-nil, nil" or "nil, err".
func (m *Manager) LoadConfig() (*Config, error) {
	// 1. get rawConfig from config.yaml file
	rawConfig, err := getRawConfigFromFile(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// 2. get all globals vars
	globalVars, err := rawConfig.GetGlobalVars()
	if err != nil {
		return nil, err
	}

	// 3. get Tools from Apps
	appTools, err := rawConfig.GetToolsFromApps(globalVars)
	if err != nil {
		return nil, err
	}

	// 4. get Tools from rawConfig
	tools, err := rawConfig.GetTools(globalVars)
	if err != nil {
		return nil, err
	}

	tools = append(tools, appTools...)

	config := &Config{
		PluginDir: rawConfig.PluginDir,
		State:     rawConfig.State,
		Tools:     tools,
	}

	// 5. validate config
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}
