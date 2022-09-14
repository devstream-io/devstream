package helm

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

type InstanceState struct {
	Workflows Workflows
}

func (is *InstanceState) ToStringInterfaceMap() (map[string]interface{}, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	defer encoder.Close()
	encoder.SetIndent(2)
	err := encoder.Encode(&is.Workflows)
	if err != nil {
		return nil, err
	}
	wfs := buf.String()

	return map[string]interface{}{
		"workflows": wfs,
	}, nil
}

type Workflows struct {
	Deployments  []Deployment  `yaml:"deployments,omitempty"`
	Daemonsets   []Daemonset   `yaml:"daemonsets,omitempty"`
	Statefulsets []Statefulset `yaml:"statefulsets,omitempty"`
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
	Name  string `yaml:"name"`
	Ready bool   `yaml:"ready"`
}

type Daemonset struct {
	Name  string `yaml:"name"`
	Ready bool   `yaml:"ready"`
}

type Statefulset struct {
	Name  string `yaml:"name"`
	Ready bool   `yaml:"ready"`
}
