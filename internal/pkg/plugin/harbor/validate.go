package harbor

import (
	helmCommon "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/pkg/util/helm"
)

// validate validates the options provided by the core.
func validate(options *helmCommon.Options) []error {
	return helm.DefaultsAndValidate(options.GetHelmParam())
}
