package configloader

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/pkg/util/log"
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
	} else {
		log.Debugf("Got Backend from config: %s", c.State.Backend)
		if c.State.Backend == "local" {
			if c.State.Options.StateFile == "" {
				log.Debugf("The stateFile has not been set, default value %s will be used.", local.DefaultStateFile)
				c.State.Options.StateFile = local.DefaultStateFile
			}
		} else if c.State.Backend == "s3" {
			if c.State.Options.Bucket == "" {
				errors = append(errors, fmt.Errorf("state s3 Bucket is empty"))
			}
			if c.State.Options.Region == "" {
				errors = append(errors, fmt.Errorf("state s3 Region is empty"))
			}
			if c.State.Options.Key == "" {
				errors = append(errors, fmt.Errorf("state s3 Key is empty"))
			}
		} else {
			errors = append(errors, fmt.Errorf("backend type error"))
		}
	}

	return errors
}
