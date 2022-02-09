package golang

import (
	"bytes"
	"html/template"

	"github.com/mitchellh/mapstructure"

	github "github.com/merico-dev/stream/pkg/util/github"
)

func renderTemplate(workflow *github.Workflow, options *Options) (string, error) {
	var jobs Jobs
	err := mapstructure.Decode(options.Jobs, &jobs)
	if err != nil {
		return "", err
	}

	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("githubactions").Delims("[[", "]]").Parse(workflow.WorkflowContent)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, jobs)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
