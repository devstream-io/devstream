package helm

import "github.com/devstream-io/devstream/pkg/util/validator"

func Validate(param *HelmParam) []error {
	return validator.Struct(param)
}
