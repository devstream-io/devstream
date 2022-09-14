package helm

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates helm param
func Validate(param *HelmParam) []error {
	var retErrs = validator.Struct(param)
	if param.Chart.ChartPath == "" && (param.Repo.Name == "" || param.Repo.URL == "") {
		err := fmt.Errorf("if chartPath == \"\", then the repo.Name & repo.URL must be set")
		retErrs = append(retErrs, err)
	}
	return retErrs
}
