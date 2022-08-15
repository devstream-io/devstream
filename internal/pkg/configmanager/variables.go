package configmanager

import (
	"github.com/devstream-io/devstream/pkg/util/template"
)

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	str, err := template.New().
		FromContent(fileContent).
		AddProcessor(template.AddDotForVariablesInConfigProcessor()).
		DefaultRender(fileContent, variables).
		Render()

	if err != nil {
		return nil, err
	}

	return []byte(str), nil
}
