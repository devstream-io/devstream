package generic

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
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
	ciFileContent, err := template.New().
		FromContent(ciTemplateString).
		SetDefaultRender("ci", opts.TemplateVariables).Render()

	if err != nil {
		return nil, fmt.Errorf("execute template error: %s", err)
	}

	// update file
	client, err := opts.newGitlabClient()
	if err != nil {
		return nil, err
	}
	commitInfo := &git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap: git.GitFileContentMap{
			ciFileName: []byte(ciFileContent),
		},
	}
	// the only difference between create and update
	if err = client.UpdateFiles(commitInfo); err != nil {
		return nil, err
	}

	return buildState(&opts), nil
}
