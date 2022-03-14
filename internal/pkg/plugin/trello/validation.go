package trello

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

	return retErrors
}
