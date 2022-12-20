package jiragithub

import (
	"github.com/mitchellh/mapstructure"

	github "github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func renderTemplate(workflow *github.Workflow, options *Options) (string, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return "", err
	}

	return template.NewRenderClient(&template.TemplateOption{
		Name: "jiragithub",
	}, template.ContentGetter).Render(workflow.WorkflowContent, opts)
}
