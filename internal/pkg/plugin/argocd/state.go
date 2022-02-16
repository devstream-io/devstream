package argocd

var DefaultDeploymentList = []string{
	"argocd-application-controller",
	"argocd-dex-server",
	"argocd-redis",
	"argocd-repo-server",
	"argocd-server",
}

type InstanceState struct {
	Workflows Workflows
}

func (is *InstanceState) ToStringInterfaceMap() map[string]interface{} {
	return map[string]interface{}{
		"workflows": is.Workflows,
	}
}

type Workflows struct {
	Deployments []Deployment
}

func (w *Workflows) AddDeployment(name string, ready bool) {
	if w.Deployments == nil {
		w.Deployments = make([]Deployment, 0)
	}
	w.Deployments = append(w.Deployments, Deployment{
		Name:  name,
		Ready: ready,
	})
}

type Deployment struct {
	Name  string
	Ready bool
}

func GetStaticState() *InstanceState {
	retState := &InstanceState{}
	for _, dpName := range DefaultDeploymentList {
		retState.Workflows.AddDeployment(dpName, true)
	}
	return retState
}
