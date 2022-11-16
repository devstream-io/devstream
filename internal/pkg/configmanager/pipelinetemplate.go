package configmanager

import (
	"fmt"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

type (
	PipelineTemplate struct {
		Name    string     `yaml:"name"`
		Type    string     `yaml:"type"`
		Options RawOptions `yaml:"options"`
	}
)

func findTemplateByName(name string, templates []PipelineTemplate) *PipelineTemplate {
	for _, template := range templates {
		if template.Name == name {
			return &template
		}
	}
	return nil
}

func renderCICDFromPipeTemplates(cicds []CICD, templates []PipelineTemplate, globalVars map[string]any,
	appIndex int, appName string, ciOrCD string) ([]PipelineTemplate, error) {
	cicdsOneApp := make([]PipelineTemplate, 0)
	for j, cicd := range cicds {
		switch cicd.Type {
		case "template":
			if cicd.TemplateName == "" {
				return nil, fmt.Errorf("apps[%d](%s).%s[%d].templateName is required", appIndex, appName, ciOrCD, j)
			}
			template := findTemplateByName(cicd.TemplateName, templates)
			if template == nil {
				return nil, fmt.Errorf("apps[%d](%s).%s[%d].%s not found in pipelineTemplates", appIndex, appName, ciOrCD, j, cicd.TemplateName)
			}
			// merge options
			// we shoould first merge options here
			// then render vars(globalVars, cicd.vars) in the options or in the whole template
			template = template.DeepCopy()
			if err := mergo.Merge(&template.Options, cicd.Options, mergo.WithOverride); err != nil {
				return nil, err
			}
			templateBytes, err := yaml.Marshal(template)
			if err != nil {
				return nil, err
			}

			allVars := mergeMaps(globalVars, cicd.Vars)

			renderedTemplateBytes, err := renderConfigWithVariables(string(templateBytes), allVars)
			if err != nil {
				return nil, err
			}
			templateNew := PipelineTemplate{}
			err = yaml.Unmarshal(renderedTemplateBytes, &templateNew)
			if err != nil {
				return nil, err
			}

			cicdsOneApp = append(cicdsOneApp, templateNew)
		default:
			// that means it's a plugin
			templateNew := PipelineTemplate{
				Name:    cicd.Type,
				Type:    cicd.Type,
				Options: cicd.Options,
			}
			cicdsOneApp = append(cicdsOneApp, templateNew)
		}
	}
	return cicdsOneApp, nil
}

func (template *PipelineTemplate) DeepCopy() *PipelineTemplate {
	templateBytes, err := yaml.Marshal(template)
	if err != nil {
		return nil
	}
	templateNew := PipelineTemplate{}
	err = yaml.Unmarshal(templateBytes, &templateNew)
	if err != nil {
		return nil
	}
	return &templateNew
}
