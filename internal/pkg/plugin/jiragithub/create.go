package jiragithub

import (
	"fmt"

	"github.com/merico-dev/stream/pkg/util/github"
)

// Create sets up jira-github-integ workflows.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	content, err := renderTemplate(workflow, opts)
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

	return BuildState(opts.Owner, opts.Repo), nil
}

func BuildState(owner, repo string) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", owner, repo)
	return res
}
