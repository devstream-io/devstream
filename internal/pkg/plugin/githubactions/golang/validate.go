package golang

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/github"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := github.NewGithubActionOptions(options)
	if err != nil {
		return nil, err
	}
	errs := validateGolangStruct(opts)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}
	return options, nil
}

// validate validates the options provided by the core.
func validateGolangStruct(opts *github.GithubActionOptions) []error {
	retErrors := make([]error, 0)
	if errs := validator.Struct(opts); len(errs) != 0 {
		retErrors = append(retErrors, errs...)
	}
	if opts.Test == nil {
		retErrors = append(retErrors, fmt.Errorf("golang project must have test file"))
	}
	// too complex to validate automatically
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
