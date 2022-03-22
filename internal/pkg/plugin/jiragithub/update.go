package jiragithub

import (
	"github.com/merico-dev/stream/pkg/util/github"
)

// Update remove and set up jira-github-integ workflows.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	opts, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	ghOptions := &github.Option{
		Owner:    opts.Owner,
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

	content, err := renderTemplate(workflow, opts)
	if err != nil {
		return nil, err
	}
	workflow.WorkflowContent = content

	if err = ghClient.AddWorkflow(workflow, opts.Branch); err != nil {
		return nil, err
	}

	return BuildState(opts.Owner, opts.Repo), nil
}
