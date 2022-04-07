package generic

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
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

	client, err := gitlab.NewClient()
	if err != nil {
		return nil, err
	}

	exists, err := client.FileExists(opts.PathWithNamespace, opts.Branch, ciFileName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return buildState(&opts), nil
}
