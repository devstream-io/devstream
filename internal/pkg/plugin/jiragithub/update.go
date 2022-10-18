package jiragithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

// Update remove and set up jira-github-integ workflows.
func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("options are illegal")
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

	if err := ghClient.DeleteWorkflow(workflow, opts.Branch); err != nil {
		return nil, err
	}

	content, err := renderTemplate(workflow, &opts)
	if err != nil {
		return nil, err
	}
	workflow.WorkflowContent = content

	if err = ghClient.AddWorkflow(workflow, opts.Branch); err != nil {
		return nil, err
	}

	return BuildStatus(opts.Owner, opts.Repo), nil
}
