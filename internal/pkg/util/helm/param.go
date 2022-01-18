package helm

// HelmParam is the struct for parameters with helm style.
type HelmParam struct {
	Repo  Repo
	Chart Chart
}

// Repo is the struct containing details of a git repository.
// TODO(daniel-hutao): make the Repo equals to repo.Entry
type Repo struct {
	Name string
	URL  string
}

// Chart is the struct containing details of a helm chart.
// TODO(daniel-hutao): make the Chart equals to helmclient.ChartSpec
type Chart struct {
	Name            string
	ReleaseName     string `mapstructure:"release_name"`
	Namespace       string
	CreateNamespace bool `mapstructure:"create_namespace"`
}
