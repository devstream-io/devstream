package golang

import "fmt"

type RegistryType string

const (
	RegistryDockerhub RegistryType = "dockerhub"
	RegistryHarbor    RegistryType = "harbor"
)

// Build is the struct for githubactions job.
type Build struct {
	Enable  bool
	Command string
}

// Test is the struct for githubactions job.
type Test struct {
	Enable   bool
	Command  string
	Coverage Coverage
}

// Docker is the struct for githubactions job.
type Docker struct {
	Enable   bool
	Registry Registry
}

type Registry struct {
	// only support dockerhub now
	Type       RegistryType
	Username   string
	Repository string
}

// Coverage is the struct for githubactions job.
type Coverage struct {
	Enable  bool
	Profile string
	Output  string
}

// Tag is the struct for githubactions job.
type Tag struct {
}

// Image is the struct for githubactions job.
type Image struct {
}

func (d *Docker) Validate() []error {
	retErrors := make([]error, 0)

	if !d.Enable {
		return retErrors
	}

	if d.Registry.Username == "" {
		retErrors = append(retErrors, fmt.Errorf("registry username is empty"))
	}

	// default to dockerhub
	if d.Registry.Type == "" {
		d.Registry.Type = RegistryDockerhub
		return retErrors
	}
	if d.Registry.Type != RegistryDockerhub && d.Registry.Type != RegistryHarbor {
		retErrors = append(retErrors, fmt.Errorf("registry type != (dockerhub || harbor)"))
	}

	return retErrors
}
