package pluginmanager

import (
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type (
	status  = string
	toolKey = string

	// traverser is a plugin manager that traverses the tools in patch by dependency
	// each patch returns a list of tools whose dependencies are all visited
	// you can think of this as a simplified version of finding all nodes with entry degree 0 of a directed acyclic graph
	traverser struct {
		traverseMap map[toolKey]status
		tools       []configloader.Tool
	}
)

const (
	unvisited      status = ""
	defaultVisited status = "visited"
)

// NewTraverser creates a new traverser which traverses the plugins in patch by dependency
func NewTraverser(tools []configloader.Tool) *traverser {
	return &traverser{
		traverseMap: make(map[toolKey]status),
		tools:       tools,
	}
}

// Traverse traverses the tools in patch by dependency
// it will terminate the traversal as soon as it encounters an error
// the status of the tool being traversed can be set at traversal time via SetStatus
// if the status is not set manually, it will be set to the default status "visited"
func (t *traverser) Traverse(op func(tool *configloader.Tool) error) error {
	for {

		// find the next tools patch to visit
		tools := t.nextPatch()
		if len(tools) == 0 {
			break
		}

		for _, tool := range tools {
			if err := op(&tool); err != nil {
				return err
			}

			// The one who calls this function should set the status of the tool being traversed
			// if the status is not set manually, it will be set to the default status "visited"
			if t.traverseMap[tool.Key()] == unvisited {
				t.traverseMap[tool.Key()] = defaultVisited
			}
		}
	}
	return nil
}

// GoTraverse likes Traverse, but it traverses parallel
func (t *traverser) GoTraverse(op func(tool *configloader.Tool) error) error {
	for {

		// find the next tools patch to visit
		tools := t.nextPatch()
		if len(tools) == 0 {
			break
		}

		var eg errgroup.Group
		for _, tool := range tools {
			// define a new var to avoid each loop use the same tool
			tool := tool
			// set up a goroutine for each tool
			eg.Go(func() error {
				if err := op(&tool); err != nil {
					return err
				}

				// The one who calls this function should set the status of the tool being traversed
				// if the status is not set manually, it will be set to the default status "visited"
				if t.traverseMap[tool.Key()] == unvisited {
					t.traverseMap[tool.Key()] = defaultVisited
				}
				return nil
			})
		}
		// any goroutine encounters an error, the program immediately returns
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}

// nextPatch returns the next patch of the tools which are not visited and whose dependencies are all visited
func (t *traverser) nextPatch() []configloader.Tool {
	var nextPatch []configloader.Tool
	for _, tool := range t.tools {
		// skip the tool that has been visited
		if t.IfToolVisited(&tool) {
			continue
		}

		// select all tools whose dependencies are all visited
		if t.IfDependenciesOfToolAllVisited(&tool) {
			nextPatch = append(nextPatch, tool)
		}
	}

	patchKeys := make([]toolKey, 0, len(nextPatch))
	for _, tool := range nextPatch {
		patchKeys = append(patchKeys, tool.Key())
	}
	log.Debugf("Next patch: %v", strings.Join(patchKeys, ", "))

	return nextPatch

}

// SetStatus sets the specific traverse status of the tool
func (t *traverser) SetStatus(tool *configloader.Tool, status status) {
	t.traverseMap[tool.Key()] = status
}

// GetStatus gets the specific traverse status of the tool
func (t *traverser) GetStatus(tool *configloader.Tool) status {
	return t.traverseMap[tool.Key()]
}

// IfToolVisited returns true if the tool has been visited
func (t *traverser) IfToolVisited(tool *configloader.Tool) bool {
	return t.traverseMap[tool.Key()] != unvisited
}

// IfDependenciesOfToolAllVisited returns true if all dependencies of the tool have been visited
func (t *traverser) IfDependenciesOfToolAllVisited(tool *configloader.Tool) bool {
	for _, dep := range tool.TrimmedDependsOn() {
		if t.traverseMap[dep] == unvisited {
			return false
		}
	}
	return true
}

// IfDependsOnStatusAnyMatch  returns true if any dependencies of the tool match the status
func (t *traverser) IfDependsOnStatusAnyMatch(tool *configloader.Tool, status status) bool {
	for _, dep := range tool.TrimmedDependsOn() {
		if t.traverseMap[dep] == status {
			return true
		}
	}
	return false
}

// IfDependsOnStatusAllMatch  returns true if all dependencies of the tool match the status
func (t *traverser) IfDependsOnStatusAllMatch(tool *configloader.Tool, status status) bool {
	for _, dep := range tool.TrimmedDependsOn() {
		if t.traverseMap[dep] != status {
			return false
		}
	}
	return true
}

// CheckAllVisited Check if all tools have been visited
// if Traverse does not return error, this function will return true unless code of traverser is wrong
func (t *traverser) CheckAllVisited() bool {
	if len(t.traverseMap) != len(t.tools) {
		return false
	}

	traverseCount := 0
	for _, tool := range t.tools {
		if t.traverseMap[tool.Key()] != unvisited {
			traverseCount++
		}
	}
	if traverseCount != len(t.tools) {
		return false
	}

	return true
}

// GetTools returns all the tools set by the traverser
func (t *traverser) GetTools() []configloader.Tool {
	return t.tools
}
