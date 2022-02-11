package golang

func (b *Build) Validate() []error {
	retErrors := make([]error, 0)

	// TODO(daniel-hutao): what should we validate here?

	return retErrors
}

func (t *Test) Validate() []error {
	retErrors := make([]error, 0)

	// TODO(daniel-hutao): what should we validate here?

	return retErrors
}

func (d *Docker) Validate() []error {
	retErrors := make([]error, 0)

	// TODO(daniel-hutao): what should we validate here?

	return retErrors
}

// Build is the struct for githubacions job.
type Build struct {
	Enable  bool
	Command string
}

// Test is the struct for githubacions job.
type Test struct {
	Enable   bool
	Command  string
	Coverage Coverage
}

// Docker is the struct for githubacions job.
type Docker struct {
	Enable bool
	Repo   string
}

// Coverage is the struct for githubacions job.
type Coverage struct {
	Enable  bool
	Profile string
	Output  string
}

// Tag is the struct for githubacions job.
type Tag struct {
}

// Image is the struct for githubacions job.
type Image struct {
}
