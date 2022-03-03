package jiragithub

import (
	"github.com/merico-dev/stream/pkg/util/github"
)

// Delete remove jira-github-integ workflows.
func Delete(options map[string]interface{}) (bool, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return false, err
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}

	if err := ghClient.DeleteWorkflow(workflow, opt.Branch); err != nil {
		return false, err
	}

	return true, nil
}
