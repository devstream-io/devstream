package argocd

// Param is the struct for parameters used by the argocd package.
type Param struct {
	Repo  Repo
	Chart Chart
}

// Repo is the struct containing details of a git repository.
type Repo struct {
	Name string
	URL  string
}

// Chart is the struct containing details of a helm chart.
type Chart struct {
	Name            string
	ReleaseName     string `mapstructure:"release_name"`
	Namespace       string
	CreateNamespace bool `mapstructure:"create_namespace"`
}
