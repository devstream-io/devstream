package golang

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(opts *Options) []error {
	retErrors := make([]error, 0)

	if errs := validator.Struct(opts); len(errs) != 0 {
		retErrors = append(retErrors, errs...)
	}
	// too complex to validate automatically
	if opts.Docker == nil {
		return retErrors
	}

	if errs := opts.Docker.Validate(); len(errs) != 0 {
		for _, e := range errs {
			retErrors = append(retErrors, fmt.Errorf("docker is invalid: %s", e))
		}
	}

	return retErrors
}
