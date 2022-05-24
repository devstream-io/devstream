package python

import (
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(opts *Options) []error {
	return validator.Struct(opts)
}
