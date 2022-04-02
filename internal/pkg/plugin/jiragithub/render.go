package jiragithub

import (
	"bytes"
	"html/template"

	"github.com/mitchellh/mapstructure"

	github "github.com/devstream-io/devstream/pkg/util/github"
)

func renderTemplate(workflow *github.Workflow, options *Options) (string, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return "", err
	}

	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("jiragithub").Delims("[[", "]]").Parse(workflow.WorkflowContent)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, opts)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
