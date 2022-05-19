package python

import (
	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
)

// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string       `validate:"required_without=Org"`
	Org      string       `validate:"required_without=Owner"`
	Repo     string       `validate:"required"`
	Branch   string       `validate:"required"`
	Language *ga.Language `validate:"required"`
}
