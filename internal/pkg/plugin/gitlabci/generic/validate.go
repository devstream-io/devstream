package generic

import "fmt"

// validate validates the options provided by the core.
func validate(options *Options) []error {
	retErrors := make([]error, 0)

	if options.PathWithNamespace == "" {
		retErrors = append(retErrors, fmt.Errorf("pathWithNamespace is empty"))
	}
	if options.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}
	if options.TemplateURL == "" {
		retErrors = append(retErrors, fmt.Errorf("templateURL is empty"))
	}

	return retErrors
}
