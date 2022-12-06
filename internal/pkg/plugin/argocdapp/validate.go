package argocdapp

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}
	return options, nil
}
