package helm

import (
	"bytes"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInstanceState_ToStringInterfaceMap(t *testing.T) {
	wf := Workflows{
		Deployments:  []Deployment{{"Deployment", true}, {"Deployment2", false}},
		Daemonsets:   []Daemonset{{"Daemonset", true}},
		Statefulsets: []Statefulset{{"Statefulset", true}},
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	defer encoder.Close()
	encoder.SetIndent(2)
	err := encoder.Encode(&wf)
	if err != nil {
		t.Error(err)
	}

	wfs := buf.String()

	tests := []struct {
		name      string
		Workflows Workflows
		want      map[string]interface{}
	}{
		// TODO: Add test cases.
		{"base", wf, map[string]interface{}{"workflows": wfs}},
		// {"base encode error", nil, map[string]interface{}{"workflows": struct{}{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := &InstanceState{
				Workflows: tt.Workflows,
			}
			if got := is.ToStringInterfaceMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceState.ToStringInterfaceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflows_AddDeployment(t *testing.T) {
	one, two := Deployment{"Deployment", true}, Deployment{"Deployment2", false}
	want, want2 := []Deployment{one}, []Deployment{one, two}
	wf := Workflows{
		Deployments:  []Deployment{},
		Daemonsets:   []Daemonset{},
		Statefulsets: []Statefulset{},
	}
	wf2 := Workflows{
		Deployments:  []Deployment{one},
		Daemonsets:   []Daemonset{},
		Statefulsets: []Statefulset{},
	}
	var wf3 Workflows
	tests := []struct {
		name    string
		wfs     Workflows
		element Deployment
		want    []Deployment
	}{
		// TODO: Add test cases.
		{"base empty", wf, one, want},
		{"base not empty", wf2, two, want2},
		{"base nil", wf3, one, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wfs.AddDeployment(tt.element.Name, tt.element.Ready)
			got, want := tt.wfs.Deployments, tt.want
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Workflows.AddDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflows_AddDaemonset(t *testing.T) {
	one, two := Daemonset{"Daemonset", true}, Daemonset{"Daemonset2", false}
	want, want2 := []Daemonset{one}, []Daemonset{one, two}
	wf := Workflows{
		Deployments:  []Deployment{},
		Daemonsets:   []Daemonset{},
		Statefulsets: []Statefulset{},
	}
	wf2 := Workflows{
		Deployments:  []Deployment{},
		Daemonsets:   []Daemonset{one},
		Statefulsets: []Statefulset{},
	}
	var wf3 Workflows
	tests := []struct {
		name    string
		wfs     Workflows
		element Daemonset
		want    []Daemonset
	}{
		// TODO: Add test cases.
		{"base empty", wf, one, want},
		{"base not empty", wf2, two, want2},
		{"base nil", wf3, one, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wfs.AddDaemonset(tt.element.Name, tt.element.Ready)
			got, want := tt.wfs.Daemonsets, tt.want
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Workflows.AddDaemonSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflows_AddStatefulset(t *testing.T) {
	one, two := Statefulset{"Statefulset", true}, Statefulset{"Statefulset", false}
	want, want2 := []Statefulset{one}, []Statefulset{one, two}
	wf := Workflows{
		Deployments:  []Deployment{},
		Daemonsets:   []Daemonset{},
		Statefulsets: []Statefulset{},
	}
	wf2 := Workflows{
		Deployments:  []Deployment{},
		Daemonsets:   []Daemonset{},
		Statefulsets: []Statefulset{one},
	}
	var wf3 Workflows
	tests := []struct {
		name    string
		wfs     Workflows
		element Statefulset
		want    []Statefulset
	}{
		// TODO: Add test cases.
		{"base empty", wf, one, want},
		{"base not empty", wf2, two, want2},
		{"base nil", wf3, one, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wfs.AddStatefulset(tt.element.Name, tt.element.Ready)
			got, want := tt.wfs.Statefulsets, tt.want
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Workflows.AddDaemonSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
