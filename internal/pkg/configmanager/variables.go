package configmanager

import (
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	logFunc(fileContent, variables)

	str, err := template.NewRenderClient(&template.TemplateOption{
		Name: "configmanager",
	}, template.ContentGetter, template.AddDotForVariablesInConfigProcessor, template.AddQuoteForVariablesInConfigProcessor,
	).Render(fileContent, variables)

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
