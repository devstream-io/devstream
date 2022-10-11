package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

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

	client, err := opts.newGitlabClient()
	if err != nil {
		return nil, err
	}

	pathInfo, err := client.GetPathInfo(ciFileName)
	if err != nil {
		return nil, err
	}

	if len(pathInfo) == 0 {
		return nil, nil
	}

	return buildState(&opts), nil
}
