package jiragithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Delete remove jira-github-integ workflows.
func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("options are illegal")
	}

	ghOptions := &github.Option{
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
