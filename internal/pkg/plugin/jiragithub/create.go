package jiragithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Create sets up jira-github-integ workflows.
func Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
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

	content, err := renderTemplate(workflow, &opts)
	if err != nil {
		return nil, err
	}
	workflow.WorkflowContent = content

	if err := ghClient.AddWorkflow(workflow, opts.Branch); err != nil {
		return nil, err
	}

	if err := setRepoSecrets(ghClient); err != nil {
		return nil, err
	}

	return BuildStatus(opts.Owner, opts.Repo), nil
}

func BuildStatus(owner, repo string) statemanager.ResourceStatus {
	resStatus := make(statemanager.ResourceStatus)
	resStatus["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", owner, repo)
	return resStatus
}
