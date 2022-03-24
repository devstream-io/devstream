package golang

import "fmt"

func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	if opts.PathWithNamespace == "" {
		retErrors = append(retErrors, fmt.Errorf("pathWithNamespace is empty"))
	}

	if opts.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	return retErrors
}
