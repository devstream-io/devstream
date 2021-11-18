package argocd

type Param struct {
	Repo  Repo
	Chart Chart
}

type Repo struct {
	Name string
	URL  string
}

type Chart struct {
	Name            string
	ReleaseName     string `mapstructure:"release_name"`
	Namespace       string
	CreateNamespace bool `mapstructure:"create_namespace"`
}
