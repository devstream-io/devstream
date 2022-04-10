package trellogithub

import "fmt"

func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	// owner/org/repo/branch
	if opts.Owner == "" && opts.Org == "" {
		retErrors = append(retErrors, fmt.Errorf("owner and org are empty"))
	}
	if opts.Repo == "" {
		retErrors = append(retErrors, fmt.Errorf("repo is empty"))
	}
	if opts.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	// TODO(daniel-hutao): Add more validations after refactor this package.

	return retErrors
}
