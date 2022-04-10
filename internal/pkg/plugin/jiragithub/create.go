package jiragithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create sets up jira-github-integ workflows.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	ghOptions := &github.Option{
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

	return BuildState(opts.Owner, opts.Repo), nil
}

func BuildState(owner, repo string) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", owner, repo)
	return res
}
