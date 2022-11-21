package helminstaller

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	return helm.Validate(options)
}
