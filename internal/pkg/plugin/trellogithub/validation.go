package trellogithub

import "fmt"

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

	// TODO(daniel-hutao): Add more validations after refactor this package.

	return retErrors
}
