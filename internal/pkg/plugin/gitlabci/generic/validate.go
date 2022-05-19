package generic

import (
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options *Options) []error {
	return validator.Struct(options)
}
