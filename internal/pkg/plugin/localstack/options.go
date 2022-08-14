package localstack

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
)

// Options is the struct for configurations of the localstack plugin.
type (
	Options struct {
		helm.Options `mapstructure:",squash"`
	}
)
