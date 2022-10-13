package ci

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the dtm-core.
func Validate(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	fieldErr := validator.StructAllError(opts)
	if fieldErr != nil {
		return nil, fieldErr
	}
	return options, nil
}

// SetDefaultConfig will update options empty values base on import options
func SetDefaultConfig(defaultConfig *Options) plugininstaller.MutableOperation {
	return func(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}
		opts.FillDefaultValue(defaultConfig)
		return opts.Encode()
	}
}
