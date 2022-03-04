package configloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyPass(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: Plugin{Kind: "argocd"}},
		{Name: "argocdapp", Plugin: Plugin{Kind: "argocdapp"}, DependsOn: "argocd.argocd"},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0, "Dependency check passed.")

}

func TestDependencyNotExist(t *testing.T) {
	tools := []Tool{
		{Name: "argocdapp", Plugin: Plugin{Kind: "argocdapp"}, DependsOn: "argocd.argocd"},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 1)

}

func TestMultipleDependencies(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: Plugin{Kind: "argocd"}},
		{Name: "repo", Plugin: Plugin{Kind: "github"}},
		{Name: "argocdapp", Plugin: Plugin{Kind: "argocdapp"}, DependsOn: "argocd.argocd,repo.github"},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}

func TestEmptyDependency(t *testing.T) {
	tools := []Tool{
		{Name: "argocd", Plugin: Plugin{Kind: "argocd"}},
		{Name: "argocdapp", Plugin: Plugin{Kind: "argocdapp"}, DependsOn: ""},
	}
	errors := validateDependency(tools)
	assert.Equal(t, len(errors), 0)
}
