package configmanager

type Manager struct {
	ConfigFile string
}

func NewManager(configFileName string) *Manager {
	return &Manager{
		ConfigFile: configFileName,
	}
}

// LoadConfig reads an input file as a general config.
// It will return "non-nil, err" or "nil, err".
func (m *Manager) LoadConfig() (*Config, error) {
	// 1. get new RawConfig from ConfigFile
	RawConfig, err := newRawConfig(m.ConfigFile)
	if err != nil {
		return nil, err
	}
	// 2. get all globals vars
	globalVars, err := RawConfig.getGlobalVars()
	if err != nil {
		return nil, err
	}
	// 3. get Apps from RawConfig
	appTools, err := RawConfig.getAppsTools(globalVars)
	if err != nil {
		return nil, err
	}

	// 4. get Tools from RawConfig
	tools, err := RawConfig.getTools(globalVars)
	if err != nil {
		return nil, err
	}
	tools = append(tools, appTools...)
	config := &Config{
		PluginDir: RawConfig.PluginDir,
		State:     RawConfig.State,
		Tools:     tools,
	}
	//5. validate config is valid
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}
