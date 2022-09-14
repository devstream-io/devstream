package helm

import "github.com/devstream-io/devstream/pkg/util/types"

// HelmParam is the struct for parameters with helm style.
type HelmParam struct {
	Repo  Repo
	Chart Chart
}

// Repo is the struct containing details of a git repository.
// TODO(daniel-hutao): make the Repo equals to repo.Entry
type Repo struct {
	// if Name or URL equals to "", then Chart.ChartPath must be set
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

// Chart is the struct containing details of a helm chart.
// TODO(daniel-hutao): make the Chart equals to helmclient.ChartSpec
type Chart struct {
	// if ChartPath equals to "", then Repo.Name and Repo.URL must be set
	ChartPath   string `mapstructure:"chartPath"`
	ChartName   string `validate:"required" mapstructure:"chartName"`
	Version     string `mapstructure:"version"`
	ReleaseName string `mapstructure:"releaseName"`
	Namespace   string `mapstructure:"namespace"`
	Wait        *bool  `mapstructure:"wait"`
	Timeout     string `mapstructure:"timeout"` // such as "1.5h" or "2h45m", valid time units are "s", "m", "h"
	UpgradeCRDs *bool  `mapstructure:"upgradeCRDs"`
	// ValuesYaml is the values.yaml content.
	// use string instead of map[string]interface{}
	ValuesYaml string `validate:"omitempty,yaml" mapstructure:"valuesYaml"`
}

func (repo *Repo) FillDefaultValue(defaultRepo *Repo) {
	types.FillStructDefaultValue(repo, defaultRepo)
}

func (chart *Chart) FillDefaultValue(defaultChart *Chart) {
	types.FillStructDefaultValue(chart, defaultChart)
}
