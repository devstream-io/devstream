package configmanager

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// rawConfig is used to describe original raw configs read from files
type rawConfig struct {
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
	totalConfigBytes  []byte         `mapstructure:"-"`
}

// getRawConfigFromFile generate new rawConfig options
func getRawConfigFromFile(configFilePath string) (*rawConfig, error) {
	// 1. get baseDir from configFile
	baseDir, err := file.GetFileAbsDirPath(configFilePath)
	if err != nil {
		return nil, err
	}

	// 2. read the original main config file
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	// TODO(daniel-hutao): We should change the documents to delete all "---" with config file. After a while, delete the following line of code, or prompt the user with "This is the wrong way" when "---" be detected

	// replace all "---", otherwise yaml.Unmarshal can only read the content before the first "---"
	configBytes = bytes.Replace(configBytes, []byte("---"), []byte("\n"), -1)

	// 3. decode total yaml files to get rawConfig
	var rawConfig rawConfig
	if err := yaml.Unmarshal(configBytes, &rawConfig); err != nil && !errors.Is(err, io.EOF) {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}

	rawConfig.configFileBaseDir = baseDir
	rawConfig.totalConfigBytes = configBytes
	return &rawConfig, nil
}

// GetGlobalVars will get global variables from GlobalVars field and varFile content
func (r *rawConfig) GetGlobalVars() (map[string]any, error) {
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

// UnmarshalYAML is used for rawConfig
// it will put variables fields in globalVars field
func (r *rawConfig) UnmarshalYAML(value *yaml.Node) error {
	configMap := make(map[string]any)
	if err := value.Decode(configMap); err != nil {
		return err
	}
	return mapstructure.Decode(configMap, r)
}

// GetToolsFromApps will get Tools from rawConfig.totalConfigBytes config
func (r *rawConfig) GetToolsFromApps(globalVars map[string]any) (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.apps[*]"
	appFileBytes, err := r.AppFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}

	_, appArray, err := getMergedNodeConfig(appFileBytes, r.totalConfigBytes, yamlPath)
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

// GetTools get Tools from tool config
func (r *rawConfig) GetTools(globalVars map[string]any) (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.tools[*]"
	fileBytes, err := r.ToolFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}
	toolStr, _, err := getMergedNodeConfig(fileBytes, r.totalConfigBytes, yamlPath)
	if err != nil {
		return nil, err
	}

	// 2. render config str with global variables
	toolsWithRenderdStr, err := renderConfigWithVariables(toolStr, globalVars)
	if err != nil {
		return nil, err
	}

	//3. unmarshal config str to Tools
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
func (r *rawConfig) getPipelineTemplatesMap() (map[string]string, error) {
	yamlPath := "$.pipelineTemplates[*]"
	fileBytes, err := r.TemplateFile.getContentBytes(r.configFileBaseDir)
	if err != nil {
		return nil, err
	}
	_, templateArray, err := getMergedNodeConfig(fileBytes, r.totalConfigBytes, yamlPath)
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

// configFileLoc is configFile location
type configFileLoc string

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

	fileBytes, err := os.ReadFile(fileAbs)
	if err != nil {
		return nil, err
	}
	return fileBytes, err
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
