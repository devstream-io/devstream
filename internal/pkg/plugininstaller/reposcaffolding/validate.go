package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the core.
func Validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	err = validator.StructAllError(opts)
	if err != nil {
		return nil, err
	}
	return options, nil
}
