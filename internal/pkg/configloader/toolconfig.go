package configloader

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Tool is the struct for one section of the DevStream tool file (part of the config.)
type Tool struct {
	Name string `yaml:"name" validate:"required"`
	// RFC 1123 - DNS Subdomain Names style
	// contain no more than 253 characters
	// contain only lowercase alphanumeric characters, '-' or '.'
	// start with an alphanumeric character
	// end with an alphanumeric character
	InstanceID string                 `yaml:"instanceID" validate:"required,dns1123subdomain"`
	DependsOn  []string               `yaml:"dependsOn"`
	Options    map[string]interface{} `yaml:"options"`
}

func (t *Tool) Validate() []error {
	return validator.Struct(t)
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

func (t *Tool) Key() string {
	return fmt.Sprintf("%s.%s", t.Name, t.InstanceID)
}

func NewToolWithToolConfigFileAndVarsConfigFile(toolFilePath, varFilePath string) ([]Tool, error) {
	toolConfigBytes, err := readFile(toolFilePath)
	if err != nil {
		return nil, err
	}

	if varFilePath == "" {
		return NewToolWithToolConfigBytesAndVarsConfigBytes(toolConfigBytes, nil)
	}

	varConfigBytes, err := readFile(varFilePath)
	if err != nil {
		return nil, err
	}

	return NewToolWithToolConfigBytesAndVarsConfigBytes(toolConfigBytes, varConfigBytes)
}

func NewToolWithToolConfigBytesAndVarsConfigBytes(toolConfigBytes, varConfigBytes []byte) ([]Tool, error) {
	if len(varConfigBytes) == 0 {
		return newToolWithToolConfig(toolConfigBytes)
	}

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
		log.Errorf("Failed to render tool config with vars: %s.", err)
		return nil, err
	}

	return newToolWithToolConfig(toolConfigBytesWithVars)
}

func newToolWithToolConfig(toolConfigBytes []byte) ([]Tool, error) {
	var config = Config{Tools: make([]Tool, 0)}
	if err := yaml.Unmarshal(toolConfigBytes, &config); err != nil {
		return nil, err
	}

	return config.Tools, nil
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

	log.Debugf("Variables file: \n%s\n", string(fileBytes))
	return fileBytes, nil
}
