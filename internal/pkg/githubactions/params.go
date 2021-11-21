package githubactions

import (
	"context"

	"github.com/google/go-github/v40/github"
)

// Workflow is the struct for a GitHub Actions workflow.
type Workflow struct {
	commitMessage    string
	workflowFileName string
	workflowContent  string
}

// Param is the struct for parameters used by the githubactions plugin.
type Param struct {
	ctx      *context.Context
	client   *github.Client
	options  *Options
	workflow *Workflow
}

// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string
	Repo     string
	Language *Language
	Branch   string
}

// Language is the struct containing details of a programming language specified in the GitHub Actions workflow.
type Language struct {
	Name    string
	Version string
}
