package pluginengine

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func dependencyResolved(tool *configmanager.Tool, unprocessedNodeSet map[string]bool) bool {
	res := true

	for _, dep := range tool.DependsOn {
		// if the tool's dependency is still not processed yet / still in the graph
		log.Debugf("TOOL %s.%s dependency NOT solved\n", tool.Name, tool.InstanceID)
		if _, ok := unprocessedNodeSet[dep]; ok {
			res = false
			break
		}
	}

	log.Debugf("TOOL %s %s %t\n", tool.Name, tool.InstanceID, res)
	return res
}

func topologicalSort(tools configmanager.Tools) ([]configmanager.Tools, error) {
	// the final result that contains sorted Tools
	// it's a sorted/ordered slice,
	// each element is a slice of Tools that can run parallel without any particular order
	res := make([]configmanager.Tools, 0)

	// a "graph", which contains "nodes" that haven't been processed yet
	unprocessedNodeSet := make(map[string]bool)
	for _, tool := range tools {
		unprocessedNodeSet[tool.KeyWithNameAndInstanceID()] = true
	}

	// while there is still a node in the graph left to be processed:
	for len(unprocessedNodeSet) > 0 {
		// the next batch of tools that can run in parallel
		batch := make(configmanager.Tools, 0)

		for _, tool := range tools {
			// if the tool has already been processed (not in the unprocessedNodeSet anymore), pass
			if _, ok := unprocessedNodeSet[tool.KeyWithNameAndInstanceID()]; !ok {
				continue
			}

			// if there isn't any dependency: it's the "start" of the graph
			// we can put it into the first batch
			if len(tool.DependsOn) == 0 {
				log.Debugf("TOOL %s.%s dependency already solved\n", tool.Name, tool.InstanceID)
				batch = append(batch, tool)
			} else {
				if dependencyResolved(tool, unprocessedNodeSet) {
					log.Debugf("TOOL %s.%s dependency already solved\n", tool.Name, tool.InstanceID)
					batch = append(batch, tool)
				}
			}
		}
		log.Debugf("BATCH: %v", batch)

		// there are still nodes unprocessed but there is no node whose dependency is solved
		// this means there might be a loop in the graph
		if len(batch) == 0 && len(unprocessedNodeSet) > 0 {
			return res, fmt.Errorf("dependency loop detected in the config")
		}

		// remove tools from the unprocessedNodeSet because they have been added to the batch
		for _, tool := range batch {
			delete(unprocessedNodeSet, tool.KeyWithNameAndInstanceID())
		}

		// add the batch to the final result
		res = append(res, batch)
	}

	return res, nil
}
