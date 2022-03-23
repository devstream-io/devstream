package golang

import "fmt"

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

	return retErrors
}
