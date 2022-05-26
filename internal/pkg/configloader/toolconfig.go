package configloader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// Tool is the struct for one section of the DevStream configuration file.
type Tool struct {
	Name string `yaml:"name"`
	// RFC 1123 - DNS Subdomain Names style
	// contain no more than 253 characters
	// contain only lowercase alphanumeric characters, '-' or '.'
	// start with an alphanumeric character
	// end with an alphanumeric character
	InstanceID string                 `yaml:"instanceID"`
	DependsOn  []string               `yaml:"dependsOn"`
	Options    map[string]interface{} `yaml:"options"`
}

func NewToolWithToolConfigFileAndVarsConfigFile(toolFilePath, varFilePath string) ([]Tool, error) {
	toolConfigBytes, err := readFile(toolFilePath)
	if err != nil {
		return nil, err
	}
	varConfigBytes, err := readFile(varFilePath)
	if err != nil {
		return nil, err
	}

	return newToolWithToolConfigAndVarsConfig(toolConfigBytes, varConfigBytes)
}

func newToolWithToolConfigAndVarsConfig(toolConfigBytes, varConfigBytes []byte) ([]Tool, error) {
	variables, err := loadVarsIntoMap(varConfigBytes)
	if err != nil {
		log.Errorf("Failed to load vars into map: %s", err)
		return nil, err
	}

	// handle variables format
	toolConfigBytesWithDot := addDotForVariablesInConfig(string(toolConfigBytes))

	// render config with variables
	toolConfigBytesWithVars, err := renderConfigWithVariables(toolConfigBytesWithDot, variables)
	if err != nil {
		log.Error("Failed to render tool config with vars: %s", err)
		return nil, err
	}

	var tools = make([]Tool, 0)
	if err := yaml.Unmarshal(toolConfigBytesWithVars, &tools); err != nil {
		return nil, err
	}

	return tools, nil
}

func loadVarsIntoMap(varConfigBytes []byte) (map[string]interface{}, error) {
	variables := make(map[string]interface{})
	err := yaml.Unmarshal(varConfigBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func readFile(filePath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Failed to read the %s: %s", filePath, err)
		return nil, err
	}

	log.Debugf("Variables file: \n%s\n", string(filePath))
	return fileBytes, nil
}

func (t *Tool) DeepCopy() *Tool {
	var retTool = Tool{
		Name:       t.Name,
		InstanceID: t.InstanceID,
		DependsOn:  t.DependsOn,
		Options:    map[string]interface{}{},
	}
	for k, v := range t.Options {
		retTool.Options[k] = v
	}
	return &retTool
}

// renderToolConfigWithVarsToConfig reads tool file rendering by var file as a Config struct.
func renderToolConfigWithVarsToConfig(toolFileName, varFileName string) (*Config, error) {
	toolFileBytes, err := ioutil.ReadFile(toolFileName)
	if err != nil {
		log.Errorf("Failed to read the toolFile: %s", err)
		return nil, err
	}

	log.Debugf("Tool config: \n%s\n", string(toolFileBytes))

	// handle variables in the config file if var file is provided
	configFileBytesWithVarsRendered, err := renderVariables(varFileName, toolFileBytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debugf("tool config file after rendering with variables: \n%s\n", string(configFileBytesWithVarsRendered))

	var config Config
	err = yaml.Unmarshal(configFileBytesWithVarsRendered, &config)
	if err != nil {
		log.Error("Please verify the format of your config file.")
		log.Errorf("Reading config file failed. %s.", err)
		return nil, err
	}

	errs := validateConfig(&config)
	if len(errs) != 0 {
		var errStrings []string
		for _, e := range errs {
			log.Errorf("Config validation failed: %s.", e)
			errStrings = append(errStrings, e.Error())
		}
		return nil, fmt.Errorf(strings.Join(errStrings, "; "))
	}

	return &config, nil
}
