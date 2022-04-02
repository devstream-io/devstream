package nodejs

import (
	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
)

// TODO(daniel-hutao): Options should keep as same as other plugins named Param
// Options is the struct for configurations of the githubactions plugin.
type Options struct {
	Owner    string
	Repo     string
	Branch   string
	Language *ga.Language
}
