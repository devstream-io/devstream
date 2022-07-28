package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	client, err := gitlab.NewClient(gitlab.WithBaseURL(opts.BaseURL))
	if err != nil {
		return false, err
	}

	if err = client.DeleteSingleFile(opts.PathWithNamespace, opts.Branch, commitMessage, ciFileName); err != nil {
		return false, err
	}
	return false, nil
}
