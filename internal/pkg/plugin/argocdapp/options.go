package argocdapp

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

// Param is the struct for parameters used by the argocdapp package.
type Options struct {
	App         App
	Destination Destination
	Source      Source
}

// App is the struct for an ArgoCD app.
type App struct {
	Name      string `validate:"required,dns1123subdomain"`
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
	Path      string `validate:"required"`
	RepoURL   string `validate:"required"`
}

// / NewOptions create options by raw options
func NewOptions(options plugininstaller.RawOptions) (Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}
