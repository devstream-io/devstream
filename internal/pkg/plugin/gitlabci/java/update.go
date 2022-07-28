package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
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

	// set with default value
	if err := opts.complete(); err != nil {
		return nil, err
	}

	// generate .gitla-ci.yml file
	content, err := renderTmpl(&opts)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.NewClient(gitlab.WithBaseURL(opts.BaseURL))
	if err != nil {
		return nil, err
	}

	// the only difference between "Create" and "Update"
	if err = client.UpdateSingleFile(opts.PathWithNamespace, opts.Branch, commitMessage, ciFileName, content); err != nil {
		return nil, err
	}

	return buildState(&opts), nil
}
