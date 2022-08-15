package generic

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// download template
	ciTemplateString, err := download(opts.TemplateURL)
	if err != nil {
		return nil, err
	}

	// render template
	var ciFileContentBytes bytes.Buffer
	tpl, err := template.New("ci").Option("missingkey=error").Parse(ciTemplateString)
	if err != nil {
		return nil, fmt.Errorf("parse template error: %s", err)
	}
	err = tpl.Execute(&ciFileContentBytes, opts.TemplateVariables)
	if err != nil {
		return nil, fmt.Errorf("execute template error: %s", err)
	}

	// commit file
	client, err := opts.newGitlabClient()
	if err != nil {
		return nil, err
	}
	_, err = client.PushLocalFileToRepo(&git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap: git.GitFileContentMap{
			ciFileName: ciFileContentBytes.Bytes(),
		},
	})
	if err != nil {
		return nil, err
	}

	return buildState(&opts), nil
}
