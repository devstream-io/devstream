package configmanager

import (
	"fmt"
	"os"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
)

type fileType string

const (
	varFile      fileType = "varfile"
	appFile      fileType = "appfile"
	toolFile     fileType = "toolfile"
	templateFile fileType = "templateFile"
)

// rawConfig is used to describe original raw configs read from files
type rawConfig struct {
	VarFile           string       `yaml:"varFile" mapstructure:"varFile"`
	ToolFile          string       `yaml:"toolFile" mapstructure:"toolFile"`
	AppFile           string       `yaml:"appFile" mapstructure:"appFile"`
	TemplateFile      string       `yaml:"templateFile" mapstructure:"templateFile"`
	State             *State       `yaml:"state" mapstructure:"state"`
	Apps              []RawOptions `yaml:"apps"`
	Tools             []RawOptions `yaml:"tools"`
	PipelineTemplates []RawOptions `yaml:"pipelineTemplates"`

	GlobalVars        map[string]any `yaml:"-" mapstructure:",remain"`
	configFileBaseDir string         `mapstructure:"-"`
	totalConfigBytes  []byte         `mapstructure:"-"`
}

// getAllTools will return the tools transfer from apps and out of apps.
func (r *rawConfig) getAllTools() (Tools, error) {
	// 1. get all globals vars
	err := r.mergeGlobalVars()
	if err != nil {
		return nil, err
	}

	// 2. get Tools from Apps
	appTools, err := r.getToolsFromApps()
	if err != nil {
		return nil, err
	}

	// 3. get Tools out of apps
	tools, err := r.getToolsOutOfApps()
	if err != nil {
		return nil, err
	}

	tools = append(tools, appTools...)
	return tools, nil
}

// mergeGlobalVars will merge the global vars from varFile and rawConfig, then merge it to rawConfig.GlobalVars
func (r *rawConfig) mergeGlobalVars() error {
	valueContent, err := r.getFileBytes(varFile)
	if err != nil {
		return err
	}

	globalVars := make(map[string]any)
	if err = yaml.Unmarshal(valueContent, globalVars); err != nil {
		return err
	}

	if err = mergo.Merge(&(r.GlobalVars), globalVars); err != nil {
		return err
	}
	return nil
}

// getToolsFromApps will get Tools from rawConfig.totalConfigBytes config
func (r *rawConfig) getToolsFromApps() (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.apps[*]"
	appFileBytes, err := r.getFileBytes(appFile)
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
		appTools, err := getToolsFromApp(appConfigStr, r.GlobalVars, templateMap)
		if err != nil {
			return nil, err
		}
		tools = append(tools, appTools...)
	}
	return tools, nil
}

// getToolsOutOfApps get Tools from tool config
func (r *rawConfig) getToolsOutOfApps() (Tools, error) {
	// 1. get tools config str
	yamlPath := "$.tools[*]"
	fileBytes, err := r.getFileBytes(toolFile)
	if err != nil {
		return nil, err
	}
	toolStr, _, err := getMergedNodeConfig(fileBytes, r.totalConfigBytes, yamlPath)
	if err != nil {
		return nil, err
	}

	// 2. render config str with global variables
	toolsWithRenderdStr, err := renderConfigWithVariables(toolStr, r.GlobalVars)
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
	fileBytes, err := r.getFileBytes(templateFile)
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

func (r *rawConfig) getFileBytes(fType fileType) ([]byte, error) {
	switch fType {
	case varFile:
		return getFileBytesInSpecificDir(r.VarFile, r.configFileBaseDir)
	case appFile:
		return getFileBytesInSpecificDir(r.AppFile, r.configFileBaseDir)
	case toolFile:
		return getFileBytesInSpecificDir(r.ToolFile, r.configFileBaseDir)
	case templateFile:
		return getFileBytesInSpecificDir(r.TemplateFile, r.configFileBaseDir)
	default:
		return nil, fmt.Errorf("not likely to happen")
	}
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

// getFileBytesInSpecificDir get file content with abs path for configFile
func getFileBytesInSpecificDir(filePath, baseDir string) ([]byte, error) {
	// if configFile is not setted, return empty content
	if filePath == "" {
		return []byte{}, nil
	}

	// refer other config file path by directory of main config file
	fileAbsPath, err := file.GenerateAbsFilePath(baseDir, filePath)
	if err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(fileAbsPath)
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
