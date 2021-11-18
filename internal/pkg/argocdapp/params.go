package argocdapp

type Param struct {
	App         App
	Destination Destination
	Source      Source
}

type App struct {
	Name      string
	Namespace string
}

type Destination struct {
	Server    string
	Namespace string
}

type Source struct {
	Valuefile string
	Path      string
	RepoURL   string
}
