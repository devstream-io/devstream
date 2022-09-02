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

	return template.New().FromContent(workflow.WorkflowContent).SetDefaultRender("jiragithub", opts).Render()
}
