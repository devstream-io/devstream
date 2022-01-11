package argocd

type Argocd struct {
	name string
}

func NewArgocd() *Argocd {
	return &Argocd{name: "ArgoCD"}
}

func (a *Argocd) Health() (bool, error) {
	return true, nil
}

func (a *Argocd) Name() string {
	return a.name
}
