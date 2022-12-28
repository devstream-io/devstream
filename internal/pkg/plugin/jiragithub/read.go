package jiragithub

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Read get jira-github-integ workflows.
func Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	if err := validator.CheckStructError(&opts).Combine(); err != nil {
		return nil, err
	}

	ghOptions := &git.RepoInfo{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	path, err := ghClient.GetWorkflowPath()
	if err != nil {
		return nil, err
	}

	resStatus := make(statemanager.ResourceStatus)
	resStatus["workflowDir"] = path
	return resStatus, nil
}
