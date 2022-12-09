package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the core.
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	if err := validator.StructAllError(opts); err != nil {
		return nil, err
	}
	return options, nil
}
