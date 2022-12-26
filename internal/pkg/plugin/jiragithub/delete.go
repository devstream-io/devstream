package jiragithub

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Delete remove jira-github-integ workflows.
func Delete(options configmanager.RawOptions) (bool, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return false, err
	}

	if err := validator.CheckStructError(&opts).Combine(); err != nil {
		return false, err
	}

	ghOptions := &git.RepoInfo{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}

	if err := ghClient.DeleteWorkflow(workflow, opts.Branch); err != nil {
		return false, err
	}

	return true, nil
}
