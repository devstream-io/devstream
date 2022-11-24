package configmanager

import (
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	log.Debugf("renderConfigWithVariables got str: %s", fileContent)
	log.Debugf("Vars: ---")
	for k, v := range variables {
		log.Debugf("%s: %s", k, v)
	}
	log.Debugf("---")

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
