package golang

import (
	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
)

// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string
	Org      string
	Repo     string
	Branch   string
	Language *ga.Language
	Build    *Build
	Test     *Test
	Docker   *Docker
}
