package configmanager

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	GOOS   string = runtime.GOOS
	GOARCH string = runtime.GOARCH
)

// Config records rendered config values and is used as a general config in DevStream.
type Config struct {
	PluginDir string
	Tools     []Tool `yaml:"tools"`
	State     *State
}

func (c *Config) Validate() []error {
	errors := make([]error, 0)

	for _, t := range c.Tools {
		errors = append(errors, t.Validate()...)
	}

	errors = append(errors, c.ValidateDependency()...)

	return errors
}

func (c *Config) ValidateDependency() []error {
	errors := make([]error, 0)

	// config "set" (map)
	toolMap := make(map[string]bool)
	// creating the set
	for _, tool := range c.Tools {
		toolMap[tool.Key()] = true
	}

	for _, tool := range c.Tools {
		// no dependency, pass
		if len(tool.DependsOn) == 0 {
			continue
		}

		// for each dependency
		for _, dependency := range tool.DependsOn {
			// skip empty string
			dependency = strings.TrimSpace(dependency)
			if dependency == "" {
				continue
			}

			// generate an error if the dependency isn't in the config set,
			if _, ok := toolMap[dependency]; !ok {
				errors = append(errors, fmt.Errorf("tool %s's dependency %s doesn't exist in the config", tool.InstanceID, dependency))
			}
		}
	}

	return errors
}

type Manager struct {
	ConfigFile string
}

func NewManager(configFileName string) *Manager {
	return &Manager{
		ConfigFile: configFileName,
	}
}

// LoadConfig reads an input file as a general config.
func (m *Manager) LoadConfig() (*Config, error) {
	// 1. read the original config file
	originalConfigFileBytes, err := m.loadOriginalConfigFile()
	if err != nil {
		return nil, err
	}

	// 2. split original config
	coreConfigBytes, variablesConfigBytes, toolsConfigBytes, err := m.splitConfigFileBytes(originalConfigFileBytes)
	if err != nil {
		return nil, err
	}

	return m.renderConfigs(coreConfigBytes, variablesConfigBytes, toolsConfigBytes)
}

func (m *Manager) renderConfigs(coreConfigBytes, variablesConfigBytes, toolsConfigBytes []byte) (*Config, error) {
	// 1. unmarshal core config file
	var coreConfig CoreConfig
	if err := yaml.Unmarshal(coreConfigBytes, &coreConfig); err != nil {
		log.Errorf("Please verify the format of your core config. Error: %s.", err)
		return nil, err
	}
	if err := coreConfig.Validate(); err != nil {
		return nil, err
	}
	state := coreConfig.State

	// 2. unmarshal tool config with variables(if needed).
	tools, err := m.renderToolsFromCoreConfigAndConfigBytes(&coreConfig, toolsConfigBytes, variablesConfigBytes)
	if err != nil {
		return nil, err
	}

	return &Config{
		PluginDir: coreConfig.PluginDir,
		Tools:     tools,
		State:     state,
	}, nil
}

func (m *Manager) renderToolsFromCoreConfigAndConfigBytes(coreConfig *CoreConfig, toolsConfigBytes, variablesConfigBytes []byte) ([]Tool, error) {
	if coreConfig.ToolFile == "" && len(toolsConfigBytes) == 0 {
		return nil, fmt.Errorf("tools config is empty")
	}

	var tools []Tool
	var err error
	if coreConfig.ToolFile == "" {
		tools, err = NewToolWithToolConfigBytesAndVarsConfigBytes(toolsConfigBytes, variablesConfigBytes)
		if err != nil {
			return nil, err
		}
	} else {
		if err = coreConfig.ParseToolFilePath(); err != nil {
			return nil, err
		}
		if coreConfig.VarFile != "" {
			if err = coreConfig.ParseVarFilePath(); err != nil {
				return nil, err
			}
		}

		tools, err = NewToolWithToolConfigFileAndVarsConfigFile(coreConfig.ToolFile, coreConfig.VarFile)
		if err != nil {
			return nil, err
		}
	}

	return tools, nil
}

func (m *Manager) loadOriginalConfigFile() ([]byte, error) {
	originalConfigFileBytes, err := os.ReadFile(m.ConfigFile)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info(`Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the "-f" option?`)
		log.Info(`See "dtm help" for more information."`)
		return nil, err
	}
	log.Debugf("Original config: \n%s\n", string(originalConfigFileBytes))
	return originalConfigFileBytes, err
}

// splitConfigFileBytes take the original config file and split it to:
// 1. core config
// 2. variable config
// 3. tool config
// Original config should be like below:
// ---
// # core config
// varFile: "" # If not empty, use the specified external variables config file
// toolFile: "" # If not empty, use the specified external tools config file
// pluginDir: "" # If empty, use the default value: ~/.devstream/plugins, or use -d flag to specify a directory
// state:
//
//	backend: local
//	options:
//	  stateFile: devstream.state
//
// ---
// # variables config
// foo: bar
//
// ---
// # plugins config
// tools:
//   - name: A-PLUGIN-NAME
//     instanceID: default
//     options:
//     foo: bar
//
// See https://github.com/devstream-io/devstream/issues/596 for more details.
func (m *Manager) splitConfigFileBytes(fileBytes []byte) (coreConfig []byte, varConfig []byte, toolConfig []byte, err error) {
	splitBytes := bytes.Split(bytes.TrimPrefix(fileBytes, []byte("---")), []byte("---"))

	if len(splitBytes) > 3 {
		err = fmt.Errorf("invalid config format")
		return
	}

	var result bool

	for _, configBytes := range splitBytes {
		if result, err = m.checkConfigType(configBytes, "core"); result {
			if len(coreConfig) > 0 {
				err = fmt.Errorf("exist multiple sections of core config")
				return
			}
			coreConfig = configBytes
			continue
		}
		if result, err = m.checkConfigType(configBytes, "tool"); result {
			if len(toolConfig) > 0 {
				err = fmt.Errorf("exist multiple sections of tool config")
				return
			}
			toolConfig = configBytes
			continue
		}
		if err != nil {
			return
		}
		if len(varConfig) > 0 {
			err = fmt.Errorf("exist multiple sections of var config")
			return
		}
		varConfig = configBytes
	}

	if len(coreConfig) == 0 {
		err = fmt.Errorf("core config is empty")
	}

	return
}

// checkConfigType checks the bytes of the configType
// core config is the core configType and can be identified by key state
// plugins config is the tool configType and can be identified by key tool
func (m *Manager) checkConfigType(bytes []byte, configType string) (bool, error) {
	result := make(map[string]interface{})
	if err := yaml.Unmarshal(bytes, &result); err != nil {
		log.Debugf("Config type checked: %s", configType)
		log.Errorf("Please verify the format of your config file. Error: %s.", err)
		return false, err
	}
	switch configType {
	case "core":
		if _, ok := result["state"]; ok {
			return true, nil
		}
	case "tool":
		if _, ok := result["tools"]; ok {
			return true, nil
		}
	}
	return false, nil
}

// GetPluginName return plugin name without file extensions
func GetPluginName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s", t.Name, GOOS, GOARCH, version.Version)
}
func GetPluginNameWithOSAndArch(t *Tool, os, arch string) string {
	return fmt.Sprintf("%s-%s-%s_%s", t.Name, os, arch, version.Version)
}

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return GetPluginName(t) + ".so"
}
func GetPluginFileNameWithOSAndArch(t *Tool, os, arch string) string {
	return GetPluginNameWithOSAndArch(t, os, arch) + ".so"
}

// GetPluginMD5FileName  If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.md5"
func GetPluginMD5FileName(t *Tool) string {
	return GetPluginName(t) + ".md5"
}
func GetPluginMD5FileNameWithOSAndArch(t *Tool, os, arch string) string {
	return GetPluginNameWithOSAndArch(t, os, arch) + ".md5"
}
