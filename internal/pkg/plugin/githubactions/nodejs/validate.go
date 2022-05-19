package nodejs

import (
	"github.com/devstream-io/devstream/pkg/util/validator"
)

func validate(opts *Options) []error {
	return validator.Struct(opts)
}
