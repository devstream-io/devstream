package golang

import (
	"fmt"

	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/mitchellh/mapstructure"
)

// validate validates the options provided by the core.
func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	// owner/repo/branch
	if opts.Owner == "" {
		retErrors = append(retErrors, fmt.Errorf("owner is empty"))
	}
	if opts.Repo == "" {
		retErrors = append(retErrors, fmt.Errorf("repo is empty"))
	}
	if opts.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	// language
	if opts.Language == nil {
		retErrors = append(retErrors, fmt.Errorf("language is empty"))
	}
	if errs := opts.Language.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("language is invalid: %s", e))
		}
	}

	// jobs
	if opts.Test == nil {
		retErrors = append(retErrors, fmt.Errorf("test is empty"))
	}
	if errs := opts.Test.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("test is invalid: %s", e))
		}
	}

	if opts.Docker == nil {
		return retErrors
	}

	if errs := opts.Docker.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("docker is invalid: %s", e))
		}
	}

	return retErrors
}

func parseAndValidateOptions(options map[string]interface{}) (*Options, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	return &opts, nil
}
