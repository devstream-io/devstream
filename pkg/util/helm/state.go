package helm

type InstanceState struct {
	Workflows Workflows
}

func (is *InstanceState) ToStringInterfaceMap() map[string]interface{} {
	return map[string]interface{}{
		"workflows": is.Workflows,
	}
}

type Workflows struct {
	Deployments  []Deployment
	Daemonsets   []Daemonset
	Statefulsets []Statefulset
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

func (w *Workflows) AddDaemonset(name string, ready bool) {
	if w.Daemonsets == nil {
		w.Daemonsets = make([]Daemonset, 0)
	}
	w.Daemonsets = append(w.Daemonsets, Daemonset{
		Name:  name,
		Ready: ready,
	})
}

func (w *Workflows) AddStatefulset(name string, ready bool) {
	if w.Statefulsets == nil {
		w.Statefulsets = make([]Statefulset, 0)
	}
	w.Statefulsets = append(w.Statefulsets, Statefulset{
		Name:  name,
		Ready: ready,
	})
}

type Deployment struct {
	Name  string
	Ready bool
}

type Daemonset struct {
	Name  string
	Ready bool
}

type Statefulset struct {
	Name  string
	Ready bool
}
