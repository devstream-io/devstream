package golang

import "github.com/merico-dev/stream/internal/pkg/util/github"

// Install installs github-repo-scaffolding-golang with provided options.
func Install(options *map[string]interface{}) (bool, error) {
	// TODO(daniel-hutao): implement it
	_, _ = github.NewGithubClient()
	_ = validate()

	return true, nil
}
