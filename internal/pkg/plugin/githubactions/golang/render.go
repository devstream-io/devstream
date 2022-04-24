package golang

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
	// let Build is false when empty
	if opts.Build == nil {
		opts.Build = &Build{false, ""}
	}

	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("githubactions").Delims("[[", "]]").Parse(workflow.WorkflowContent)
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
