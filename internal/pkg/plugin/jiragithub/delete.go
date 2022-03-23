package jiragithub

import (
	"github.com/merico-dev/stream/pkg/util/github"
)

// Delete remove jira-github-integ workflows.
func Delete(options map[string]interface{}) (bool, error) {
	opts, err := parseAndValidateOptions(options)
	if err != nil {
		return false, err
	}

	ghOptions := &github.Option{
		Owner:    opts.Owner,
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
