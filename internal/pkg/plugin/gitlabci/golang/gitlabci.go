package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
)

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

type Options struct {
	PathWithNamespace string
	Branch            string
}

func parseAndValidateOptions(options *map[string]interface{}) (*Options, error) {
	var opt Options
	err := mapstructure.Decode(options, &opt)
	if err != nil {
		return nil, err
	}

	if errs := validateParameters(&opt); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return nil, fmt.Errorf("incorrect params")
	}

	return &opt, nil
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

func buildState(opt *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opt.PathWithNamespace,
		"branch":            opt.Branch,
	}
}
