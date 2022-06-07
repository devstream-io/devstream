package configloader

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

// Config is the struct for loading DevStream configuration YAML files.
// It records rendered config values and is used as a general config in DevStream.
type Config struct {
	Tools []Tool `yaml:"tools"`
	State *State
}

// LoadConfig reads an input file as a general config.
func LoadConfig(configFileName string) (*Config, error) {
	// 1. read the original config file
	originalConfigFileBytes, err := loadOriginalConfigFile(configFileName)
	if err != nil {
		return nil, err
	}

	// 2. split original config
	coreConfigBytes, variablesConfigBytes, toolsConfigBytes, err := SplitConfigFileBytes(originalConfigFileBytes)
	if err != nil {
		return nil, err
	}

	return renderConfigs(coreConfigBytes, variablesConfigBytes, toolsConfigBytes)
}

func renderConfigs(coreConfigBytes, variablesConfigBytes, toolsConfigBytes []byte) (*Config, error) {
	// 1. unmarshal core config file
	var coreConfig CoreConfig
	if err := yaml.Unmarshal(coreConfigBytes, &coreConfig); err != nil {
		log.Errorf("Please verify the format of your core config. Error: %s.", err)
		return nil, err
	}
	if ok, err := coreConfig.Validate(); !ok {
		return nil, err
	}
	state := coreConfig.State

	// 2. unmarshal tool config with variables(if needed).
	tools, err := renderToolsFromCoreConfigAndConfigBytes(&coreConfig, toolsConfigBytes, variablesConfigBytes)
	if err != nil {
		return nil, err
	}

	return &Config{
		Tools: tools,
		State: state,
	}, nil
}

func renderToolsFromCoreConfigAndConfigBytes(coreConfig *CoreConfig, toolsConfigBytes, variablesConfigBytes []byte) ([]Tool, error) {
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

func loadOriginalConfigFile(configFile string) ([]byte, error) {
	originalConfigFileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info("Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil, err
	}
	log.Debugf("Original config: \n%s\n", string(originalConfigFileBytes))
	return originalConfigFileBytes, err
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

// SplitConfigFileBytes take the original config file and split it to:
// 1. core config
// 2. variables config
// 3. tools config
// Original config should be like below:
// ---
// # core config
// varFile: "" # If not empty, use the specified external variables config file
// toolFile: "" # If not empty, use the specified external tools config file
// state:
//   backend: local
//   options:
//     stateFile: devstream.state
//
// ---
// # variables config
// foo: bar
//
// ---
// # plugins config
// tools:
// - name: A-PLUGIN-NAME
//   instanceID: default
//   options:
//     foo: bar
//
// See https://github.com/devstream-io/devstream/issues/596 for more details.
func SplitConfigFileBytes(fileBytes []byte) (coreConfig []byte, varConfig []byte, toolConfig []byte, err error) {
	splitBytes := bytes.Split(bytes.TrimPrefix(fileBytes, []byte("---")), []byte("---"))

	switch len(splitBytes) {
	case 1:
		coreConfig = splitBytes[0]
	case 2:
		coreConfig = splitBytes[0]
		toolConfig = splitBytes[1]
	case 3:
		coreConfig = splitBytes[0]
		varConfig = splitBytes[1]
		toolConfig = splitBytes[2]
	default:
		err = fmt.Errorf("invalid config format")
	}

	if len(coreConfig) == 0 {
		err = fmt.Errorf("core config is empty")
	}

	return
}

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.so", t.Name, GOOS, GOARCH, version.Version)
}

// GetPluginMD5FileName  If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.md5"
func GetPluginMD5FileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.md5", t.Name, GOOS, GOARCH, version.Version)
}
