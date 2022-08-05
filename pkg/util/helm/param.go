package helm

import (
	"github.com/devstream-io/devstream/pkg/util/types"
)

// HelmParam is the struct for parameters with helm style.
type HelmParam struct {
	Repo  Repo
	Chart Chart
}

// Repo is the struct containing details of a git repository.
// TODO(daniel-hutao): make the Repo equals to repo.Entry
type Repo struct {
	Name string `validate:"required" mapstructure:"name"`
	URL  string `validate:"required" mapstructure:"url"`
}

// Chart is the struct containing details of a helm chart.
// TODO(daniel-hutao): make the Chart equals to helmclient.ChartSpec
type Chart struct {
	ChartName       string `validate:"required" mapstructure:"chart_name"`
	Version         string `mapstructure:"version"`
	ReleaseName     string `mapstructure:"release_name"`
	Namespace       string `mapstructure:"namespace"`
	CreateNamespace *bool  `mapstructure:"create_namespace"`
	Wait            *bool  `mapstructure:"wait"`
	Timeout         string `mapstructure:"timeout"` // such as "1.5h" or "2h45m", valid time units are "s", "m", "h"
	UpgradeCRDs     *bool  `mapstructure:"upgradeCRDs"`
	// ValuesYaml is the values.yaml content.
	// use string instead of map[string]interface{}
	ValuesYaml string `mapstructure:"values_yaml"`
}

func (repo *Repo) FillDefaultValue(defaultRepo *Repo) {
	if repo.Name == "" {
		repo.Name = defaultRepo.Name
	}
	if repo.URL == "" {
		repo.URL = defaultRepo.URL
	}
}

func (chart *Chart) FillDefaultValue(defaultChart *Chart) {
	if chart.ChartName == "" {
		chart.ChartName = defaultChart.ChartName
	}
	if chart.Timeout == "" {
		chart.Timeout = defaultChart.Timeout
	}
	chart.UpgradeCRDs = getBoolValue(chart.UpgradeCRDs, defaultChart.UpgradeCRDs)
	chart.Wait = getBoolValue(chart.Wait, defaultChart.Wait)
	chart.CreateNamespace = getBoolValue(chart.CreateNamespace, defaultChart.CreateNamespace)
}

func getBoolValue(field, defaultField *bool) *bool {
	if field != nil {
		return field
	}

	if defaultField != nil {
		return defaultField
	}

	return types.Bool(false)
}
