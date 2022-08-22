package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	// validate parameters
	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	client, err := opts.newGitlabClient()
	if err != nil {
		return nil, err
	}

	ciFileContent, err := client.GetGitLabCIGolangTemplate()
	if err != nil {
		return nil, err
	}

	_, err = client.PushLocalFileToRepo(&git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap: git.GitFileContentMap{
			ciFileName: []byte(ciFileContent),
		},
	}, false)
	if err != nil {
		return nil, err
	}

	return buildState(&opts), nil
}
