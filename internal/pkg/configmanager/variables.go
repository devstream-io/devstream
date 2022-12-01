package configmanager

import (
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	logFunc(fileContent, variables)

	str, err := template.New().
		FromContent(fileContent).
		AddDotForVariablesInConfigProcessor().
		AddQuoteForVariablesInConfigProcessor().
		SetDefaultRender(fileContent, variables).
		Render()

	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}

func logFunc(fileContent string, variables map[string]interface{}) {
	log.Debugf("renderConfigWithVariables got str: %s", fileContent)
	log.Debug("Vars: ---")
	for k, v := range variables {
		log.Debugf("%s: %s", k, v)
	}
	log.Debug("---")
}
