package argocdapp

// Param is the struct for parameters used by the argocdapp package.
type Options struct {
	App         App
	Destination Destination
	Source      Source
}

// App is the struct for an ArgoCD app.
type App struct {
	Name      string
	Namespace string
}

// Destination is the struct for the destination of an ArgoCD app.
type Destination struct {
	Server    string
	Namespace string
}

// Source is the struct for the source of an ArgoCD app.
type Source struct {
	Valuefile string
	Path      string
	RepoURL   string
}
