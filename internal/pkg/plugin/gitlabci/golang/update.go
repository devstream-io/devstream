package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Update(options map[string]interface{}) (map[string]interface{}, error) {
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

	client, err := gitlab.NewClient()
	if err != nil {
		return nil, err
	}

	ciFileContent, err := client.GetGitLabCIGolangTemplate()
	if err != nil {
		return nil, err
	}

	// the only difference between "Create" and "Update"
	if err = client.UpdateSingleFile(opts.PathWithNamespace, opts.Branch, commitMessage, ciFileName, ciFileContent); err != nil {
		return nil, err
	}

	return buildState(&opts), nil
}
