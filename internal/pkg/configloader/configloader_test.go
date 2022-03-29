package configloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyPass(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: "argocd"},
		{Name: "argocdapp", Plugin: "argocdapp", DependsOn: []string{"argocd.argocd"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0, "Dependency check passed.")

}

func TestDependencyNotExist(t *testing.T) {
	tools := []Tool{
		{Name: "argocdapp", Plugin: "argocdapp", DependsOn: []string{"argocd.argocd"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 1)

}

func TestMultipleDependencies(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: "argocd"},
		{Name: "repo", Plugin: "github"},
		{Name: "argocdapp", Plugin: "argocdapp", DependsOn: []string{"argocd.argocd", "repo.github"}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}

func TestEmptyDependency(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: "argocd"},
		{Name: "argocdapp", Plugin: "argocdapp", DependsOn: []string{}},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}
