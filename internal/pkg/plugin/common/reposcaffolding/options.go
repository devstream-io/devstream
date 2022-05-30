package reposcaffolding

import (
	"fmt"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

type Options struct {
	Owner             string `validate:"required_without=Org"`
	Org               string `validate:"required_without=Owner"`
	Repo              string `validate:"required"`
	Branch            string `validate:"required"`
	PathWithNamespace string
	BaseURL           string `validate:"omitempty,url"`
	ImageRepo         string `mapstructure:"image_repo"`
	Visibility        string `validate:"omitempty,oneof=public private internal"`
}

// Validate validates the options provided by the core.
func Validate(opts *Options) []error {
	retErrors := make([]error, 0)
	retErrors = append(retErrors, validator.Struct(opts)...)
	// set PathWithNamespace for GitLab. GitHub won't need to use this
	// opts.PathWithNamespace = fmt.Sprintf("%s/%s", opts.Owner, opts.Repo)
	if opts.Org != "" {
		opts.PathWithNamespace = filepath.Clean(fmt.Sprintf("%s/%s", opts.Org, opts.Repo))
	} else {
		opts.PathWithNamespace = filepath.Clean(fmt.Sprintf("%s/%s", opts.Owner, opts.Repo))
	}

	return retErrors
}
