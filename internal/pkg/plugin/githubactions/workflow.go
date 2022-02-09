package githubactions

import "fmt"

const (
	CommitMessage       = "GitHub Actions workflow, created by DevStream"
	PRBuilderFileName   = "pr-builder.yml"
	MainBuilderFileName = "main-builder.yml"
)

// Language is the struct containing details of a programming language specified in the GitHub Actions Workflow.
type Language struct {
	Name    string
	Version string
}

func (l *Language) Validate() []error {
	retErrors := make([]error, 0)

	if l.Name == "" {
		retErrors = append(retErrors, fmt.Errorf("name is empty"))
	}

	return retErrors
}
