package githubactions

import (
	"context"

	"github.com/google/go-github/v40/github"
)

type Pipeline struct {
	commitMessage    string
	workflowFileName string
	workflowContent  string
}

type Param struct {
	ctx      *context.Context
	client   *github.Client
	options  *Options
	pipeline *Pipeline
}

type Options struct {
	Owner    string
	Repo     string
	Language *Language
	Branch   string
}

type Language struct {
	Name    string
	Version string
}
