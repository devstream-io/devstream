package helm

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates helm param
func Validate(param *HelmParam) []error {
	var retErrs = validator.Struct(param)

	if param.Chart.ChartPath == "" && (param.Repo.Name == "" || param.Repo.URL == "" || param.Chart.ChartName == "") {
		log.Debugf("Repo.Name: %s, Repo.URL: %s, Chart.ChartName: %s", param.Repo.Name, param.Repo.URL, param.Chart.ChartName)
		err := fmt.Errorf("if chartPath == \"\", then the repo.Name & repo.URL & chart.chartName must be set")
		retErrs = append(retErrs, err)
	}

	return retErrs
}
