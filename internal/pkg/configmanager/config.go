package configmanager

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/imdario/mergo"
)

type (
	// configFileLoc is configFile location
	configFileLoc string
	// Config records rendered config values and is used as a general config in DevStream.
	Config struct {
		// Command line flag have a higher priority than the config file.
		// If you used the `--plugin-dir` flag with `dtm`, then the "pluginDir" in the config file will be ignored.
		PluginDir string
		Tools     Tools `yaml:"tools"`
		State     *State
	}
	// RawConfig is used to describe original raw configs read from files
	RawConfig struct {
		VarFile           configFileLoc `yaml:"varFile" mapstructure:"varFile"`
		ToolFile          configFileLoc `yaml:"toolFile" mapstructure:"toolFile"`
		AppFile           configFileLoc `yaml:"appFile" mapstructure:"appFile"`
		TemplateFile      configFileLoc `yaml:"templateFile" mapstructure:"templateFile"`
		PluginDir         string        `yaml:"pluginDir" mapstructure:"pluginDir"`
		State             *State        `yaml:"state" mapstructure:"state"`
		Apps              []RawOptions  `yaml:"apps"`
		Tools             []RawOptions  `yaml:"tools"`
		PipelineTemplates []RawOptions  `yaml:"pipelineTemplates"`

		GlobalVars        map[string]any `yaml:"-" mapstructure:",remain"`
		configFileBaseDir string         `mapstructure:"-"`
		globalBytes       []byte         `mapstructure:"-"`
	}
	// State is the struct for reading the state configuration in the config file.
	// It defines how the state is stored, specifies the type of backend and related options.
	State struct {
		Backend string             `yaml:"backend"`
		Options StateConfigOptions `yaml:"options"`
	}
	// StateConfigOptions is the struct for reading the options of the state backend.
	StateConfigOptions struct {
		// for s3 backend
		Bucket string `yaml:"bucket"`
		Region string `yaml:"region"`
		Key    string `yaml:"key"`
		// for local backend
		StateFile string `yaml:"stateFile"`
		// for ConfigMap backend
		Namespace string `yaml:"namespace"`
		ConfigMap string `yaml:"configmap"`
	}
)

func (c *Config) validate() error {
	if c.State == nil {
		return fmt.Errorf("state is not defined")
	}
	return nil
}

// newRawConfig generate new RawConfig options
func newRawConfig(configFileLocation string) (*RawConfig, error) {
	// 1. get baseDir from configFile
	baseDir, err := file.GetFileAbsDirPath(configFileLocation)
	if err != nil {
		return nil, err
	}

	// 2. read the original main config file
	configBytes, err := loadConfigFile(configFileLocation)
	if err != nil {
		return nil, err
	}
	// replace all "---"
	// otherwise yaml.Unmarshal can only read the content before the first "---"
	configBytes = bytes.Replace(configBytes, []byte("---"), []byte("\n"), -1)

	// 3. decode total yaml files to get RawConfig
	var RawConfig RawConfig
	if err := yaml.Unmarshal(configBytes, &RawConfig); err != nil && !errors.Is(err, io.EOF) {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}
	RawConfig.configFileBaseDir = baseDir
	RawConfig.globalBytes = configBytes
	return &RawConfig, nil
}

// loadConfigFile get config file content by location
func loadConfigFile(fileLoc string) ([]byte, error) {
	configBytes, err := os.ReadFile(fileLoc)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info(`Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the "-f" option?`)
		log.Info(`See "dtm help" for more information."`)
		return nil, err
	}
	log.Debugf("Original config: \n%s\n", string(configBytes))
	return configBytes, err
}

// getGlobalVars will get global variables from GlobalVars field and varFile content
func (r *RawConfig) getGlobalVars() (map[string]any, error) {
	valueContent, err := r.VarFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}
	globalVars := make(map[string]any)
	if err := yaml.Unmarshal(valueContent, globalVars); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&globalVars, r.GlobalVars); err != nil {
		return nil, err
	}
	return globalVars, nil
}

// UnmarshalYAML is used for RawConfig
// it will put variables fields in globalVars field
func (r *RawConfig) UnmarshalYAML(value *yaml.Node) error {
	configMap := make(map[string]any)
	if err := value.Decode(configMap); err != nil {
		return err
	}
	return mapstructure.Decode(configMap, r)
}

// getApps will get Apps from appStr config
func (r *RawConfig) getAppsTools(globalVars map[string]any) (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.apps[*]"
	fileBytes, err := r.AppFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}

	_, appArray, err := getMergedNodeConfig(fileBytes, r.globalBytes, yamlPath)
	if err != nil {
		return nil, err
	}
	// 2. get pipelineTemplates config map
	templateMap, err := r.getPipelineTemplatesMap()
	if err != nil {
		return nil, err
	}
	// 3. render app with pipelineTemplate and globalVars
	var tools Tools
	for _, appConfigStr := range appArray {
		appTools, err := getToolsFromApp(appConfigStr, globalVars, templateMap)
		if err != nil {
			return nil, err
		}
		tools = append(tools, appTools...)
	}
	return tools, nil
}

// getTools get Tools from tool config
func (r *RawConfig) getTools(globalVars map[string]any) (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.tools[*]"
	fileBytes, err := r.ToolFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}
	toolStr, _, err := getMergedNodeConfig(fileBytes, r.globalBytes, yamlPath)
	if err != nil {
		return nil, err
	}
	// 2. render config str with global variables
	toolsWithRenderdStr, err := renderConfigWithVariables(toolStr, globalVars)
	if err != nil {
		return nil, err
	}
	//3.unmarshal config str to Tools
	var tools Tools
	if err := yaml.Unmarshal([]byte(toolsWithRenderdStr), &tools); err != nil {
		return nil, err
	}
	//4. validate tools is valid
	if err := tools.validateAll(); err != nil {
		return nil, err
	}
	return tools, nil
}

// getPipelineTemplatesMap generate template name/rawString map
func (r *RawConfig) getPipelineTemplatesMap() (map[string]string, error) {
	yamlPath := "$.pipelineTemplates[*]"
	fileBytes, err := r.TemplateFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}
	_, templateArray, err := getMergedNodeConfig(fileBytes, r.globalBytes, yamlPath)
	if err != nil {
		return nil, err
	}
	templateMap := make(map[string]string)
	for _, templateStr := range templateArray {
		templateName, err := file.GetYamlNodeStrByPath([]byte(templateStr), "$.name")
		if err != nil {
			return nil, err
		}
		templateMap[templateName] = templateStr
	}
	return templateMap, nil
}

// getContentBytes get file content with abs path for configFile
func (f configFileLoc) getContentBytes(baseDir string) ([]byte, error) {
	// if configFile is not setted, return empty content
	if string(f) == "" {
		return []byte{}, nil
	}
	// refer other config file path by directory of main config file
	fileAbs, err := file.GenerateAbsFilePath(baseDir, string(f))
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(fileAbs)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// getMergedNodeConfig will use yamlPath to config from configFile content and global content
// then merge these content
func getMergedNodeConfig(fileBytes []byte, globalBytes []byte, yamlPath string) (string, []string, error) {
	fileNode, err := file.GetYamlNodeArrayByPath(fileBytes, yamlPath)
	if err != nil {
		return "", nil, err
	}
	globalNode, err := file.GetYamlNodeArrayByPath(globalBytes, yamlPath)
	if err != nil {
		return "", nil, err
	}
	mergedNode := file.MergeYamlNode(fileNode, globalNode)
	if mergedNode == nil {
		return "", nil, err
	}
	return mergedNode.StrOrigin, mergedNode.StrArray, nil
}
