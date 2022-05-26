package pluginengine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func TestTopologicalSortChangesInBatch(t *testing.T) {
	toolA := &configloader.Tool{InstanceID: "a", Name: "a"}
	toolB := &configloader.Tool{InstanceID: "b", Name: "b"}
	toolC := &configloader.Tool{InstanceID: "c", Name: "c", DependsOn: []string{"a.a", "b.b"}}
	toolD := &configloader.Tool{InstanceID: "d", Name: "d", DependsOn: []string{"b.b"}}
	toolE := &configloader.Tool{InstanceID: "e", Name: "e", DependsOn: []string{"c.c"}}

	changes := []*Change{
		{
			Tool:       toolA,
			ActionName: statemanager.ActionCreate,
		},
		// there is no toolB
		// because in real scenarios, some tools may not need to be changed
		// that is to say, the tools extracted from changes may not contain all the tools defined in the config
		// however, as this unit test demonstrates, the function still works properly in this situation
		{
			Tool:       toolD,
			ActionName: statemanager.ActionCreate,
		},
		{
			Tool:       toolE,
			ActionName: statemanager.ActionUpdate,
		},
		// simulates a delete operation added at the end of apply
		{
			Tool:       toolC,
			ActionName: statemanager.ActionDelete,
		},
	}

	expected := [][]*Change{
		// first batch
		{
			{
				Tool:       toolA,
				ActionName: statemanager.ActionCreate,
			},
			// although D depends on B
			// but B is not in the changes list
			// so D is in the first batch
			{
				Tool:       toolD,
				ActionName: statemanager.ActionCreate,
			},
		},

		// second batch
		{

			{
				Tool:       toolC,
				ActionName: statemanager.ActionDelete,
			},
		},

		// third batch
		{
			{
				Tool:       toolE,
				ActionName: statemanager.ActionUpdate,
			},
		},
	}

	_ = toolB

	actual, err := topologicalSortChangesInBatch(changes)

	assert.NoError(t, err)

	assert.Equal(t, len(expected), len(actual))
	for i, batch := range actual {
		assert.Equal(t, len(expected[i]), len(batch))
		for j, change := range batch {
			assert.Equal(t, *expected[i][j], *change)
		}
	}
}
