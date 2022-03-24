package argocdapp

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

// validate validates the options provided by the core.

func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	if opts.App.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("app.name is empty"))
	}
	if errs := validation.IsDNS1123Subdomain(opts.App.Name); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("app.name %s is invalid: %s", opts.App.Name, e))
		}
	}

	if opts.Source.Path == "" {
		retErrors = append(retErrors, fmt.Errorf("source.path is empty"))
	}
	if opts.Source.RepoURL == "" {
		retErrors = append(retErrors, fmt.Errorf("source.repoURL is empty"))
	}

	return retErrors
}
