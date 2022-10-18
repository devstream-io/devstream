package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
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
	commitInfo := &git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap: git.GitFileContentMap{
			ciFileName: []byte(content),
		},
	}

	// the only difference between "Create" and "Update"
	if _, err = client.PushFiles(commitInfo, true); err != nil {
		return nil, err
	}

	return buildStatus(&opts), nil
}
