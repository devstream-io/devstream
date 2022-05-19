package helm

import "github.com/devstream-io/devstream/pkg/util/validator"

func Validate(param *HelmParam) []error {
	Defaults(param)
	return validator.Struct(param)
}

// Defaults set the default value with HelmParam.
// TODO(daniel-hutao): don't call this function insides the Validate()
func Defaults(param *HelmParam) {
	if param.Chart.Timeout == "" {
		// Make the timeout be same as the default value for `--timeout` with `helm install/upgrade/rollback`
		param.Chart.Timeout = "5m0s"
	}
}
