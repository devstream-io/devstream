package helm

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
)

// Validate validates the options provided by the dtm-core.
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	errs := helm.Validate(opts.GetHelmParam())
	if len(errs) > 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
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
		opts.FillDefaultValue(defaultConfig)
		return types.EncodeStruct(opts)
	}
}
