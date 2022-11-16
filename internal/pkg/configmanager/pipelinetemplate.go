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

func getTemplateByName(name string, templates []PipelineTemplate) *PipelineTemplate {
	for _, template := range templates {
		if template.Name == name {
			return &template
		}
	}
	return nil
}

func renderCICDFromPipelineTemplates(cicds []CICD, templates []PipelineTemplate, globalVars map[string]any,
	appIndex int, appName string, ciOrCD string) ([]PipelineTemplate, error) {
	cicdsOneApp := make([]PipelineTemplate, 0)
	for j, cicd := range cicds {
		switch cicd.Type {
		case "template":
			if cicd.TemplateName == "" {
				return nil, fmt.Errorf("apps[%d](%s).%s[%d].templateName is required", appIndex, appName, ciOrCD, j)
			}
			template := getTemplateByName(cicd.TemplateName, templates)
			if template == nil {
				return nil, fmt.Errorf("apps[%d](%s).%s[%d].%s not found in pipelineTemplates", appIndex, appName, ciOrCD, j, cicd.TemplateName)
			}
			// merge options
			// we shoould first merge options here
			// then render vars(globalVars, cicd.vars) in the options or in the whole templateNew
			templateNew := template.DeepCopy()
			if err := mergo.Merge(&templateNew.Options, cicd.Options, mergo.WithOverride); err != nil {
				return nil, err
			}

			// render vars to the whole template,
			// cicd.Vars is local variables, it will overwrite the keys in globalVars
			allVars := mergeMaps(globalVars, cicd.Vars)
			if err := templateNew.renderVars(allVars); err != nil {
				return nil, err
			}

			cicdsOneApp = append(cicdsOneApp, *templateNew)
		default:
			// that means it's a plugin
			templateNew := PipelineTemplate{
				Name:    cicd.Type,
				Type:    cicd.Type,
				Options: cicd.Options,
			}

			// render vars to the whole template
			if err := templateNew.renderVars(globalVars); err != nil {
				return nil, err
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

func (template *PipelineTemplate) renderVars(vars map[string]any) error {
	templateBytes, err := yaml.Marshal(template)
	if err != nil {
		return err
	}

	renderedTemplateBytes, err := renderConfigWithVariables(string(templateBytes), vars)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(renderedTemplateBytes, &template)
}
