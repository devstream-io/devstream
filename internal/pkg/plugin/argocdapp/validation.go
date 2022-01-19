package argocdapp

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

func validate(param *Param) []error {
	retErrors := make([]error, 0)

	if param.App.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("app.name is empty"))
	}
	if errs := validation.IsDNS1123Subdomain(param.App.Name); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("app.name %s is invalid: %s", param.App.Name, e))
		}
	}

	if param.Source.Path == "" {
		retErrors = append(retErrors, fmt.Errorf("source.path is empty"))
	}
	if param.Source.RepoURL == "" {
		retErrors = append(retErrors, fmt.Errorf("source.repoURL is empty"))
	}

	return retErrors
}
