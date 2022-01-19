package configloader

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
)

func (c *Config) Validate() []error {
	if len(c.Tools) == 0 {
		return []error{fmt.Errorf("config has no tools defined")}
	}

	retErrors := make([]error, 0)

	for _, t := range c.Tools {
		retErrors = append(retErrors, t.validate()...)
	}
	return retErrors
}

func (t *Tool) validate() []error {
	retErrors := make([]error, 0)

	// Name
	if t.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("name is empty"))
	}

	errs := validation.IsDNS1123Subdomain(t.Name)
	for _, e := range errs {
		retErrors = append(retErrors, fmt.Errorf("name %s invalid: %s", t.Name, e))
	}

	// Plugin
	if t.Plugin.Kind == "" {
		retErrors = append(retErrors, fmt.Errorf("plugin.kind is empty"))
	}

	return retErrors
}
