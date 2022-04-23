package configloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyPass(t *testing.T) {
	tools := []Tool{
		{InstanceID: "argocd", Name: "argocd"},
		{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0, "Dependency check passed.")

}

func TestDependencyNotExist(t *testing.T) {
	tools := []Tool{
		{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 1)

}

func TestMultipleDependencies(t *testing.T) {
	tools := []Tool{
		{InstanceID: "argocd", Name: "argocd"},
		{InstanceID: "repo", Name: "github"},
		{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd", "github.repo"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}

func TestEmptyDependency(t *testing.T) {
	tools := []Tool{
		{InstanceID: "argocd", Name: "argocd"},
		{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}
