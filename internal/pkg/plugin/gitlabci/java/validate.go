package java

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	retErrors = append(retErrors, validator.Struct(opts)...)

	if opts.Build.Enable {
		retErrors = append(retErrors, opts.Build.validate()...)
	}

	if opts.Deploy.Enable {
		retErrors = append(retErrors, opts.Deploy.validate()...)
	}

	return retErrors
}

// Validate Build struct
// 1. dockerhub username can not be empty when build is enabled
// TODO: support harbor
// 2. build target docker image name can not be empty when build is enabled
func (b *Build) validate() []error {
	retErrors := make([]error, 0)

	if b.UserName == "" {
		retErrors = append(retErrors, fmt.Errorf("dockerhub username is empty"))
	}

	if b.ImageName == "" {
		retErrors = append(retErrors, fmt.Errorf("build target docker image name is empty"))
	}

	return retErrors
}

// Validate Deploy struct
// 1. gitlab-kubernetes agent name can not be empty when deploy is enabled
func (d *Deploy) validate() []error {
	retErrors := make([]error, 0)

	if d.K8sAgentName == "" {
		retErrors = append(retErrors, fmt.Errorf("gitlab kubernetes agent name is empty"))
	}

	return retErrors
}
