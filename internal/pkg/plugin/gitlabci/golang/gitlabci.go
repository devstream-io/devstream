package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
)

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

type Options struct {
	PathWithNamespace string
	Branch            string
}

func parseAndValidateOptions(options map[string]interface{}) (*Options, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	if errs := validateParameters(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	return &opts, nil
}

func validateParameters(param *Options) []error {
	retErrors := make([]error, 0)

	if param.PathWithNamespace == "" {
		retErrors = append(retErrors, fmt.Errorf("pathWithNamespace is empty"))
	}

	if param.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	return retErrors
}

func buildState(opts *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opts.PathWithNamespace,
		"branch":            opts.Branch,
	}
}
