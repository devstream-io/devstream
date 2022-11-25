package configmanager

import (
	"fmt"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var validPipelineConfigMap = map[string]string{
	"github-actions":   "git@github.com:devstream-io/ci-template.git//github-actions",
	"jenkins-pipeline": "https://raw.githubusercontent.com/devstream-io/ci-template/main/jenkins-pipeline/general/Jenkinsfile",
	"gitlab-ci":        "",
	"argocdapp":        "",
}

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

// newPipelineTemplate will generate pipleinTemplate from pipelineRaw
func (p *pipelineRaw) newPipeline(repo *scm.SCMInfo, templateMap map[string]string, globalVars map[string]any) (*pipelineTemplate, error) {
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
			Type: p.Type,
			Name: p.Type,
		}
	}
	t.Options = p.getOptions(t.Type, repo)
	return t, nil
}

// set default value
func (p *pipelineRaw) getOptions(piplineType string, repo *scm.SCMInfo) RawOptions {
	const configLocationKey = "configLocation"
	// 1. set default configLocation
	if p.Options == nil {
		p.Options = RawOptions{}
	}
	_, configExist := p.Options[configLocationKey]
	if !configExist {
		p.Options[configLocationKey] = validPipelineConfigMap[piplineType]
	}
	// 2. extract jenkins config from options for jenkins-pipeline
	opt := RawOptions{
		"pipeline": RawOptions(p.Options),
		"scm":      RawOptions(repo.Encode()),
	}
	switch piplineType {
	case "jenkins-pipeline":
		jenkinsOpts, exist := p.Options["jenkins"]
		if exist {
			opt["jenkins"] = jenkinsOpts
		}
	}
	return opt
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

func (t *pipelineTemplate) checkError() error {
	_, valid := validPipelineConfigMap[t.Type]
	if !valid {
		return fmt.Errorf("pipeline type %s not supported for now", t.Type)
	}
	return nil
}

func (t *pipelineTemplate) getPipelineTool(appName string) (*Tool, error) {
	// 1. check pipeline type is valid
	if err := t.checkError(); err != nil {
		return nil, err
	}
	return newTool(t.Type, appName, t.Options), nil
}
