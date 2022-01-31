package helm

import "fmt"

func Validate(param *HelmParam) []error {
	Defaults(param)

	retErrors := make([]error, 0)

	if param.Repo.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("repo.name is empty"))
	}
	if param.Repo.URL == "" {
		retErrors = append(retErrors, fmt.Errorf("repo.url is empty"))
	}
	if param.Chart.ChartName == "" {
		retErrors = append(retErrors, fmt.Errorf("chart.chart_name is empty"))
	}

	return retErrors
}

// Defaults set the default value with HelmParam.
// TODO(daniel-hutao): don't call this function insides the Validate()
func Defaults(param *HelmParam) {
	if param.Chart.Timeout == "" {
		// Make the timeout be same as the default value for `--timeout` with `helm install/upgrade/rollback`
		param.Chart.Timeout = "5m0s"
	}
}
