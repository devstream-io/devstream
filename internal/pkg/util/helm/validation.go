package helm

import "fmt"

func Validate(param *HelmParam) []error {
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
