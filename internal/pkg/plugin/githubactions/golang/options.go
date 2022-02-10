package golang

import (
	"fmt"

	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
)

// TODO(daniel-hutao): Options should keep as same as other plugins named Param
// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string
	Repo     string
	Branch   string
	Language *ga.Language
	Build    *Build
	Test     *Test
	Docker   *Docker
}

// validate validates the options provided by the core.
func validate(param *Options) []error {
	retErrors := make([]error, 0)

	// owner/repo/branch
	if param.Owner == "" {
		retErrors = append(retErrors, fmt.Errorf("owner is empty"))
	}
	if param.Repo == "" {
		retErrors = append(retErrors, fmt.Errorf("repo is empty"))
	}
	if param.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	// language
	if param.Language == nil {
		retErrors = append(retErrors, fmt.Errorf("language is empty"))
	}
	if errs := param.Language.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("language is invalid: %s", e))
		}
	}

	// jobs
	if param.Test == nil {
		retErrors = append(retErrors, fmt.Errorf("test is empty"))
	}
	if errs := param.Test.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("test is invalid: %s", e))
		}
	}

	if param.Docker == nil {
		retErrors = append(retErrors, fmt.Errorf("docker is empty"))
	}
	if errs := param.Docker.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("docker is invalid: %s", e))
		}
	}

	return retErrors
}
