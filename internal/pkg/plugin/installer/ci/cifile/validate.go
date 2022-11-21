package cifile

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the dtm-core.
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
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
func SetDefaultConfig(defaultConfig *Options) installer.MutableOperation {
	return func(options configmanager.RawOptions) (configmanager.RawOptions, error) {
		opts, err := NewOptions(options)
		if err != nil {
			return nil, err
		}
		opts.fillDefaultValue(defaultConfig)
		return types.EncodeStruct(opts)
	}
}
