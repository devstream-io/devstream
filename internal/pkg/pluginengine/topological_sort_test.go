package pluginengine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func TestNoDependency(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
		{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}},
		{Name: "d", Plugin: configloader.Plugin{Kind: "d"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
				{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
				{Name: "c", Plugin: configloader.Plugin{Kind: "c"}},
				{Name: "d", Plugin: configloader.Plugin{Kind: "d"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestSingleDependency(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
		// {Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"a.a"}},
		// {Name: "d", Plugin: configloader.Plugin{Kind: "d"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
				// {Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
				// {Name: "d", Plugin: configloader.Plugin{Kind: "d"}},
			},
			{
				{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"a.a"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestMultiDependencies(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
		{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"a.a", "b.b"}},
		{Name: "d", Plugin: configloader.Plugin{Kind: "d"}, DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
				{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
			},
			{
				{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{Name: "d", Plugin: configloader.Plugin{Kind: "d"}, DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestDependencyLoop(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
		{Name: "b", Plugin: configloader.Plugin{Kind: "b"}, DependsOn: []string{"d.d"}},
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"b.b"}},
		{Name: "d", Plugin: configloader.Plugin{Kind: "d"}, DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
				{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
			},
			{
				{Name: "c", Plugin: configloader.Plugin{Kind: "c"}, DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{Name: "d", Plugin: configloader.Plugin{Kind: "d"}, DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, fmt.Errorf("dependency loop detected in the config"), err)
	assert.NotEqual(t, expectedRes, actualRes)
}
