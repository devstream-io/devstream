package configmanager

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/file"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/mapz"
)

type (
	// pipelineRaw is the raw format of app.ci/app.cd config
	pipelineRaw struct {
		Type         string     `yaml:"type" mapstructure:"type"`
		TemplateName string     `yaml:"templateName" mapstructure:"templateName"`
		Options      RawOptions `yaml:"options" mapstructure:"options"`
		Vars         RawOptions `yaml:"vars" mapstructure:"vars"`
	}
	pipelineTemplate struct {
		Name    string     `yaml:"name"`
		Type    string     `yaml:"type"`
		Options RawOptions `yaml:"options"`
	}
)

func getPipelineTemplatesMapFromConfigFile(fileBytes []byte) (map[string]string, error) {
	yamlPath := "$.pipelineTemplates[*]"
	yamlStrArray, err := file.GetYamlNodeArrayByPath(fileBytes, yamlPath)
	if err != nil {
		return nil, err
	}

	if yamlStrArray == nil {
		return make(map[string]string, 0), nil
	}

	var retMap = make(map[string]string)
	for _, tplStr := range yamlStrArray.StrArray {
		tplName, err := file.GetYamlNodeStrByPath([]byte(tplStr), "$.name")
		if err != nil {
			return nil, err
		}
		retMap[tplName] = tplStr
	}

	return retMap, nil
}

// getPipelineTemplate will generate pipleinTemplate from pipelineRaw
func (p *pipelineRaw) getPipelineTemplate(templateMap map[string]string, globalVars map[string]any) (*pipelineTemplate, error) {
	var (
		t   *pipelineTemplate
		err error
	)
	switch p.Type {
	case "template":
		t, err = p.newPipelineFromTemplate(templateMap, globalVars)
		if err != nil {
			return nil, err
		}
	default:
		t = &pipelineTemplate{
			Type:    p.Type,
			Name:    p.Type,
			Options: p.Options,
		}
	}
	return t, nil
}

func (p *pipelineRaw) newPipelineFromTemplate(templateMap map[string]string, globalVars map[string]any) (*pipelineTemplate, error) {
	var t pipelineTemplate
	if p.TemplateName == "" {
		return nil, fmt.Errorf("templateName is required")
	}
	templateStr, ok := templateMap[p.TemplateName]
	if !ok {
		return nil, fmt.Errorf("%s not found in pipelineTemplates", p.TemplateName)
	}

	allVars := mapz.Merge(globalVars, p.Vars)
	templateRenderdStr, err := renderConfigWithVariables(templateStr, allVars)
	if err != nil {
		return nil, fmt.Errorf("%s render pipelineTemplate failed: %+w", p.TemplateName, err)
	}

	if err := yaml.Unmarshal([]byte(templateRenderdStr), &t); err != nil {
		return nil, fmt.Errorf("%s parse pipelineTemplate yaml failed: %+w", p.TemplateName, err)
	}

	if err := mergo.Merge(&p.Options, t.Options); err != nil {
		return nil, fmt.Errorf("%s merge template options faield: %+v", p.TemplateName, err)
	}

	return &t, nil
}

func (t *pipelineTemplate) generatePipelineTool(app *app) (*Tool, error) {
	const configLocationKey = "configLocation"
	// 1. get configurator by template type
	pipelineConfigurator, exist := optionConfiguratorMap[t.Type]
	if !exist {
		return nil, fmt.Errorf("pipeline type [%s] not supported for now", t.Type)
	}
	// 2. set default configLocation
	if _, configLocationExist := t.Options[configLocationKey]; !configLocationExist {
		if pipelineConfigurator.hasDefaultConfig() {
			t.Options[configLocationKey] = pipelineConfigurator.defaultConfigLocation
		}
	}
	// 3. generate tool options
	pipelineFinalOptions := pipelineConfigurator.optionGeneratorFunc(t.Options, app)
	return newTool(t.Type, app.Name, pipelineFinalOptions), nil
}
