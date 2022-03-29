package pluginengine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func TestNoDependency(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: "a"},
		{Name: "b", Plugin: "b"},
		{Name: "c", Plugin: "c"},
		{Name: "d", Plugin: "d"},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: "a"},
				{Name: "b", Plugin: "b"},
				{Name: "c", Plugin: "c"},
				{Name: "d", Plugin: "d"},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestSingleDependency(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: "a"},
		{Name: "c", Plugin: "c", DependsOn: []string{"a.a"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: "a"},
			},
			{
				{Name: "c", Plugin: "c", DependsOn: []string{"a.a"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestMultiDependencies(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: "a"},
		{Name: "b", Plugin: "b"},
		{Name: "c", Plugin: "c", DependsOn: []string{"a.a", "b.b"}},
		{Name: "d", Plugin: "d", DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: "a"},
				{Name: "b", Plugin: "b"},
			},
			{
				{Name: "c", Plugin: "c", DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{Name: "d", Plugin: "d", DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestDependencyLoop(t *testing.T) {
	tools := []configloader.Tool{
		{Name: "a", Plugin: "a"},
		{Name: "b", Plugin: "b", DependsOn: []string{"d.d"}},
		{Name: "c", Plugin: "c", DependsOn: []string{"b.b"}},
		{Name: "d", Plugin: "d", DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{Name: "a", Plugin: "a"},
				{Name: "b", Plugin: "b"},
			},
			{
				{Name: "c", Plugin: "c", DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{Name: "d", Plugin: "d", DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, fmt.Errorf("dependency loop detected in the config"), err)
	assert.NotEqual(t, expectedRes, actualRes)
}
