package configloader

import (
	"fmt"
	"strings"
)

// validateDependency validates if dependency of tools are legal
// mainly used to check if all the tools defined in `dependsOn` exist
func validateDependency(tools []Tool) []error {
	errors := make([]error, 0)

	// config "set" (map)
	toolMap := make(map[string]bool)
	// creating the set
	for _, tool := range tools {
		toolMap[tool.Key()] = true
	}

	for _, tool := range tools {
		// no dependency, pass
		if len(tool.DependsOn) == 0 {
			continue
		}

		// for each dependency
		for _, dependency := range tool.DependsOn {
			// skip empty string
			dependency = strings.TrimSpace(dependency)
			if dependency == "" {
				continue
			}

			// generate an error if the dependency isn't in the config set,
			if _, ok := toolMap[dependency]; !ok {
				errors = append(errors, fmt.Errorf("tool %s's dependency %s doesn't exist in the config", tool.InstanceID, dependency))
			}
		}
	}

	return errors
}
