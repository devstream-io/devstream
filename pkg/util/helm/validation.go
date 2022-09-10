package helm

import "github.com/devstream-io/devstream/pkg/util/validator"

// Validate validates helm param
func Validate(param *HelmParam) []error {
	return validator.Struct(param)
}
