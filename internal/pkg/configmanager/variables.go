package configmanager

import (
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func getVarsFromConfigFile(fileBytes []byte) (map[string]any, error) {
	yamlPath := "$.vars"
	yamlStr, err := file.GetYamlNodeStrByPath(fileBytes, yamlPath)
	if err != nil {
		return nil, err
	}

	var retMap = make(map[string]any)
	err = yaml.Unmarshal([]byte(yamlStr), retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	if len(variables) == 0 {
		return []byte(fileContent), nil
	}
	logFunc(fileContent, variables)

	str, err := template.New().
		FromContent(fileContent).
		AddProcessor(template.AddDotForVariablesInConfigProcessor()).
		SetDefaultRender(fileContent, variables).
		Render()

	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}

func logFunc(fileContent string, variables map[string]interface{}) {
	log.Debugf("renderConfigWithVariables got str: %s", fileContent)
	log.Debugf("Vars: ---")
	for k, v := range variables {
		log.Debugf("%s: %s", k, v)
	}
	log.Debugf("---")
}
