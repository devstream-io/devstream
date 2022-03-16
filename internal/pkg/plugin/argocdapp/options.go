package argocdapp

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

// Param is the struct for parameters used by the argocdapp package.
type Options struct {
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

func validateOptions(opts *Options) []error {
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
