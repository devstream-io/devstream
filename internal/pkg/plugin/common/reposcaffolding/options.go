package reposcaffolding

import "fmt"

type Options struct {
	Owner             string
	Org               string
	Repo              string
	Branch            string
	PathWithNamespace string
	ImageRepo         string `mapstructure:"image_repo"`
}

// Validate validates the options provided by the core.
func Validate(opts *Options) []error {
	retErrors := make([]error, 0)

	// owner/org/repo/branch
	if opts.Owner == "" && opts.Org == "" {
		retErrors = append(retErrors, fmt.Errorf("owner and org are empty"))
	}
	if opts.Repo == "" {
		retErrors = append(retErrors, fmt.Errorf("repo is empty"))
	}

	// set PathWithNamespace for GitLab. GitHub won't need to use this
	opts.PathWithNamespace = fmt.Sprintf("%s/%s", opts.Owner, opts.Repo)
	if opts.Branch == "" {
		retErrors = append(retErrors, fmt.Errorf("branch is empty"))
	}

	return retErrors
}
