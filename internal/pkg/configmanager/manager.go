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
	// 1. get new ConfigRaw from ConfigFile
	configRaw, err := newConfigRaw(m.ConfigFile)
	if err != nil {
		return nil, err
	}
	// 2. get all globals vars
	globalVars, err := configRaw.getGlobalVars()
	if err != nil {
		return nil, err
	}
	// 3. get Apps from ConfigRaw
	appTools, err := configRaw.getAppsTools(globalVars)
	if err != nil {
		return nil, err
	}

	// 4. get Tools from ConfigRaw
	tools, err := configRaw.getTools(globalVars)
	if err != nil {
		return nil, err
	}
	tools = append(tools, appTools...)
	config := &Config{
		PluginDir: configRaw.PluginDir,
		State:     configRaw.State,
		Tools:     tools,
	}
	//5. validate config is valid
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}
