package helm

// HelmParam is the struct for parameters with helm style.
type HelmParam struct {
	Repo  Repo
	Chart Chart
}

// Repo is the struct containing details of a git repository.
// TODO(daniel-hutao): make the Repo equals to repo.Entry
type Repo struct {
	Name string `validate:"required"`
	URL  string `validate:"required"`
}

// Chart is the struct containing details of a helm chart.
// TODO(daniel-hutao): make the Chart equals to helmclient.ChartSpec
type Chart struct {
	ChartName       string `validate:"required" mapstructure:"chart_name"`
	Version         string
	ReleaseName     string `mapstructure:"release_name"`
	Namespace       string
	CreateNamespace bool `mapstructure:"create_namespace"`
	Wait            bool
	Timeout         string // such as "1.5h" or "2h45m", valid time units are "s", "m", "h"
	UpgradeCRDs     bool   `mapstructure:"upgradeCRDs"`
	// ValuesYaml is the values.yaml content.
	// use string instead of map[string]interface{}
	ValuesYaml string `mapstructure:"values_yaml"`
}
