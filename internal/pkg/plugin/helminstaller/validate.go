package helminstaller

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	helmOptions, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	if err := validator.CheckStructError(helmOptions).Combine(); err != nil {
		return nil, err
	}
	return options, nil
}
