package argocdapp

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

// Param is the struct for parameters used by the argocdapp package.
type Param struct {
	App         App
	Destination Destination
	Source      Source
}

// App is the struct for an ArgoCD app.
type App struct {
	Name      string
	Namespace string
}

// Destination is the struct for the destination of an ArgoCD app.
type Destination struct {
	Server    string
	Namespace string
}

// Source is the struct for the source of an ArgoCD app.
type Source struct {
	Valuefile string
	Path      string
	RepoURL   string
}

func validateParams(param *Param) []error {
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
