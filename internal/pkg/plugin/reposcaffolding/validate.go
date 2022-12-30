package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
		return nil, err
	}
	if err := validator.CheckStructError(opts).Combine(); err != nil {
		return nil, err
	}
	return rawOptions, nil
}
