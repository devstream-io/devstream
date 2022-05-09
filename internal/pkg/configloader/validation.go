package configloader

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
)

func validateConfig(config *Config) []error {
	errors := make([]error, 0)

	for _, t := range config.Tools {
		errors = append(errors, validateTool(&t)...)
	}

	errors = append(errors, validateDependency(config.Tools)...)

	return errors
}

func validateTool(t *Tool) []error {
	errors := make([]error, 0)

	// InstanceID
	if t.InstanceID == "" {
		errors = append(errors, fmt.Errorf("name is empty"))
	}

	errs := validation.IsDNS1123Subdomain(t.InstanceID)
	for _, e := range errs {
		errors = append(errors, fmt.Errorf("name %s is invalid: %s", t.InstanceID, e))
	}

	// InstanceID
	if t.Name == "" {
		errors = append(errors, fmt.Errorf("plugin is empty"))
	}

	return errors
}

func validateDependency(tools []Tool) []error {
	errors := make([]error, 0)

	// config "set" (map)
	toolMap := make(map[string]bool)
	// creating the set
	for _, tool := range tools {
		key := fmt.Sprintf("%s.%s", tool.Name, tool.InstanceID)
		toolMap[key] = true
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

// validateConfigFile validate all the general config items
func validateConfigFile(c *ConfigFile) []error {
	errors := make([]error, 0)

	if c.ToolFile == "" {
		errors = append(errors, fmt.Errorf("tool file is empty"))
	}

	if c.State == nil {
		errors = append(errors, fmt.Errorf("state config is empty"))
	}

	if c.State.Options == nil {
		errors = append(errors, fmt.Errorf("state options is empty"))
	}

	if c.State.Backend == "" {
		errors = append(errors, fmt.Errorf("backend is empty"))
	}
	return errors
}
