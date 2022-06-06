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
	var config Config
	// 1. read the original config file
	originalConfigFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info("Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil, err
	}
	log.Debugf("Original config: \n%s\n", string(originalConfigFileBytes))

	// 2. split original config
	coreConfigBytes, variablesConfigBytes, toolsConfigBytes, err := SplitConfigFileBytes(originalConfigFileBytes)
	if err != nil {
		return nil, err
	}

	if len(coreConfigBytes) == 0 {
		return nil, fmt.Errorf("core config is empty")
	}

	// 3. unmarshal core config file
	var coreConfig CoreConfig
	err = yaml.Unmarshal(coreConfigBytes, &coreConfig)
	if err != nil {
		log.Errorf("Please verify the format of your core config. Error: %s", err)
		return nil, err
	}
	config.State = coreConfig.State

	if ok, err := coreConfig.Validate(); !ok {
		return nil, err
	}

	// 4. unmarshal tool config with variables(if needed).
	if coreConfig.ToolFile == "" {
		if len(toolsConfigBytes) == 0 {
			return nil, fmt.Errorf("tools config is empty")
		}

		tools, err := NewToolWithToolConfigBytesAndVarsConfigBytes(toolsConfigBytes, variablesConfigBytes)
		if err != nil {
			return nil, err
		}
		config.Tools = tools
		return &config, nil
	}

	if err := coreConfig.ParseVarFilePath(); err != nil {
		return nil, err
	}
	if coreConfig.VarFile != "" {
		if err := coreConfig.ParseToolFilePath(); err != nil {
			return nil, err
		}
	}

	tools, err := NewToolWithToolConfigFileAndVarsConfigFile(coreConfig.ToolFile, coreConfig.VarFile)
	if err != nil {
		return nil, err
	}
	config.Tools = tools

	return &config, nil
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
// # core config (please don't change this comment)
// varFile: "" # If not empty, use the specified external variables config file
// toolFile: "" # If not empty, use the specified external tools config file
// state:
//   backend: local
//   options:
//     stateFile: devstream.state
//
// ---
// # variables config (please don't change this comment)
// foo: bar
//
// ---
// # plugins config (please don't change this comment)
// tools:
// - name: A-PLUGIN-NAME
//   instanceID: default
//   options:
//     foo: bar
//
// See https://github.com/devstream-io/devstream/issues/596 for more details.
func SplitConfigFileBytes(fileBytes []byte) ([]byte, []byte, []byte, error) {
	splitedBytes := bytes.Split(bytes.TrimPrefix(fileBytes, []byte("---")), []byte("---"))
	switch len(splitedBytes) {
	// core config only
	case 1:
		return splitedBytes[0], nil, nil, nil
		// core config + tools config
	case 2:
		return splitedBytes[0], nil, splitedBytes[2], nil
		// core config + variables config + tools config
	case 3:
		return splitedBytes[0], splitedBytes[1], splitedBytes[2], nil
	default:
		return nil, nil, nil, fmt.Errorf("invalid config format")
	}
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
