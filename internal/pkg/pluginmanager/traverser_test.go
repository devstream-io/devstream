package pluginmanager

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

func TestTraverserPatches(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"b.b", "c.c"}},
		{InstanceID: "e", Name: "e", DependsOn: []string{"d.d"}},
		{InstanceID: "f", Name: "f", DependsOn: []string{"a.a"}},
		{InstanceID: "g", Name: "g"},
	}

	traverse := NewTraverser(tools)
	patch := traverse.nextPatch()
	if !comparePatch(patch, []string{"a.a", "b.b", "g.g"}) {
		t.Errorf("Expected patch %v, got %v", []string{"a.a", "b.b", "g.g"}, toolsToKeys(patch))
	}
	for _, tool := range patch {
		traverse.SetStatus(&tool, defaultVisited)
	}

	patch = traverse.nextPatch()
	if !comparePatch(patch, []string{"c.c", "f.f"}) {
		t.Errorf("Expected patch %v, got %v", []string{"c.c", "f.f"}, toolsToKeys(patch))
	}
	for _, tool := range patch {
		traverse.SetStatus(&tool, defaultVisited)
	}

	patch = traverse.nextPatch()
	if !comparePatch(patch, []string{"d.d"}) {
		t.Errorf("Expected patch %v, got %v", []string{"d.d"}, toolsToKeys(patch))
	}
	for _, tool := range patch {
		traverse.SetStatus(&tool, defaultVisited)
	}

	patch = traverse.nextPatch()
	if !comparePatch(patch, []string{"e.e"}) {
		t.Errorf("Expected patch %v, got %v", []string{"e.e"}, toolsToKeys(patch))
	}
	for _, tool := range patch {
		traverse.SetStatus(&tool, defaultVisited)
	}

	patch = traverse.nextPatch()
	if len(patch) != 0 {
		t.Errorf("Expected empty patch, got %v", toolsToKeys(patch))
	}
}

func TestTraverserWithNoError(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"b.b", "c.c"}},
		{InstanceID: "e", Name: "e", DependsOn: []string{"d.d"}},
		{InstanceID: "f", Name: "f", DependsOn: []string{"a.a"}},
		{InstanceID: "g", Name: "g"},
	}
	traverse := NewTraverser(tools)
	traverseList := make([]string, 0)
	traverseFunc := func(tool *configloader.Tool) error {
		traverseList = append(traverseList, tool.Key())
		return nil
	}
	err := traverse.Traverse(traverseFunc)
	assert.NoError(t, err)
	assert.Equal(t, true, traverse.CheckAllVisited())
	assert.Equal(t, []string{"a.a", "b.b", "g.g", "c.c", "f.f", "d.d", "e.e"}, traverseList)
}

func TestTraverserWithError(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"b.b", "c.c"}},
		{InstanceID: "e", Name: "e", DependsOn: []string{"d.d"}},
		{InstanceID: "f", Name: "f", DependsOn: []string{"a.a"}},
		{InstanceID: "g", Name: "g"},
	}
	traverse := NewTraverser(tools)
	traverseFunc := func(tool *configloader.Tool) error {
		if tool.Key() == "c.c" {
			return errors.New("abort traversal")
		}
		return nil
	}
	err := traverse.Traverse(traverseFunc)

	// get tools that were visited
	traverseList := make([]string, 0)
	for _, tool := range tools {
		if traverse.IfToolVisited(&tool) {
			traverseList = append(traverseList, tool.Key())
		}
	}

	// because of the error, the traversal should stop at c.c
	assert.Equal(t, []string{"a.a", "b.b", "g.g"}, traverseList)

	assert.Equal(t, "abort traversal", err.Error())
}

func TestGoTraverserWithNoError(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"b.b", "c.c"}},
		{InstanceID: "e", Name: "e", DependsOn: []string{"d.d"}},
		{InstanceID: "f", Name: "f", DependsOn: []string{"a.a"}},
		{InstanceID: "g", Name: "g"},
	}
	traverse := NewTraverser(tools)
	traverseList := make([]string, 0)
	traverseFunc := func(tool *configloader.Tool) error {
		traverseList = append(traverseList, tool.Key())
		return nil
	}
	err := traverse.GoTraverse(traverseFunc)
	assert.NoError(t, err)
	assert.True(t, traverse.CheckAllVisited())

	assert.True(t, ifStringSetEqual(traverseList, []string{"a.a", "b.b", "g.g", "c.c", "f.f", "d.d", "e.e"}))
}

func TestGoTraverserWithError(t *testing.T) {
	tools := []configloader.Tool{
		{InstanceID: "a", Name: "a"},
		{InstanceID: "b", Name: "b"},
		{InstanceID: "c", Name: "c", DependsOn: []string{"a.a"}},
		{InstanceID: "d", Name: "d", DependsOn: []string{"b.b", "c.c"}},
		{InstanceID: "e", Name: "e", DependsOn: []string{"d.d"}},
		{InstanceID: "f", Name: "f", DependsOn: []string{"a.a"}},
		{InstanceID: "g", Name: "g"},
	}
	traverse := NewTraverser(tools)
	traverseFunc := func(tool *configloader.Tool) error {
		if tool.Key() == "c.c" {
			return errors.New("abort traversal")
		}
		return nil
	}
	err := traverse.GoTraverse(traverseFunc)

	// get tools that were visited
	traverseList := make([]string, 0)
	for _, tool := range tools {
		if traverse.IfToolVisited(&tool) {
			traverseList = append(traverseList, tool.Key())
		}
	}

	assert.Error(t, err, "abort traversal")

	// because of the error, the traversal should stop at c.c
	// but because of the goroutine, maybe f.f is visited, they are in same patch

	success := ifStringSetEqual(traverseList, []string{"a.a", "b.b", "g.g"}) ||
		ifStringSetEqual(traverseList, []string{"a.a", "b.b", "g.g", "f.f"})
	assert.True(t, success, fmt.Sprintf("actual tarverlist%v", traverseList))
}

// check if the patch returns is correct
func comparePatch(patch []configloader.Tool, toolKeys []string) bool {
	if len(patch) != len(toolKeys) {
		return false
	}

	pathKeys := toolsToKeys(patch)

	return ifStringSetEqual(pathKeys, toolKeys)
}

// compare two string sets if they are equal ignoring the order
func ifStringSetEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// sort two slices to compare conveniently
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// toolsToKeys returns a slice of tool keys
func toolsToKeys(tools []configloader.Tool) []string {
	keys := make([]string, 0)
	for _, tool := range tools {
		keys = append(keys, tool.Key())
	}
	return keys
}
