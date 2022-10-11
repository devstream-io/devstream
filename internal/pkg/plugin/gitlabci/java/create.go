package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
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

	// set with default value
	if err := opts.complete(); err != nil {
		return nil, err
	}

	// generate .gitla-ci.yml file
	content, err := renderTmpl(&opts)
	if err != nil {
		return nil, err
	}

	client, err := opts.newGitlabClient()
	if err != nil {
		return nil, err
	}
	_, err = client.PushFiles(&git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap: git.GitFileContentMap{
			ciFileName: []byte(content),
		},
	}, false)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
