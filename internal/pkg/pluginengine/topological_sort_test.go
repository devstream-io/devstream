package pluginengine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

func TestNoDependency(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c"},
		{InstanceID: "d", Name: "d"},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{InstanceID: "a", Name: "a"},
				{InstanceID: "b", Name: "b"},
				{InstanceID: "c", Name: "c"},
				{InstanceID: "d", Name: "d"},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestSingleDependency(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{InstanceID: "a", Name: "a"},
			},
			{
				{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestMultiDependencies(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a", "b.b"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{InstanceID: "a", Name: "a"},
				{InstanceID: "b", Name: "b"},
			},
			{
				{InstanceID: "c", Name: "c", DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{InstanceID: "d", Name: "d", DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedRes, actualRes)
}

func TestDependencyLoop(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b", DependsOn: []string{"d.d"}},
		{InstanceID: "c", Name: "c", DependsOn: []string{"b.b"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"c.c"}},
	}
	expectedRes :=
		[][]configloader.Tool{
			{
				{InstanceID: "a", Name: "a"},
				{InstanceID: "b", Name: "b"},
			},
			{
				{InstanceID: "c", Name: "c", DependsOn: []string{"a.a", "b.b"}},
			},
			{
				{InstanceID: "d", Name: "d", DependsOn: []string{"c.c"}},
			},
		}
	actualRes, err := topologicalSort(tools)
	assert.Equal(t, fmt.Errorf("dependency loop detected in the config"), err)
	assert.NotEqual(t, expectedRes, actualRes)
}
