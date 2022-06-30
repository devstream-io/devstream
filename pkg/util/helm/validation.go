package helm

import "github.com/devstream-io/devstream/pkg/util/validator"

// defaults set the default value with HelmParam.
func defaults(param *HelmParam) {
	if param.Chart.Timeout == "" {
		// Make the timeout be same as the default value for `--timeout` with `helm install/upgrade/rollback`
		param.Chart.Timeout = "5m0s"
	}
}

func validate(param *HelmParam) []error {
	return validator.Struct(param)
}

func DefaultsAndValidate(param *HelmParam) []error {
	defaults(param)
	return validate(param)
}
