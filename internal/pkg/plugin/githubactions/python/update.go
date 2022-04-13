package python

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update remove and set up GitHub Actions workflows.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
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

	log.Debugf("Language is %s.", ga.GetLanguage(opts.Language))

	for _, pipeline := range workflows {
		err := ghClient.DeleteWorkflow(pipeline, opts.Branch)
		if err != nil {
			return nil, err
		}

		err = ghClient.AddWorkflow(pipeline, opts.Branch)
		if err != nil {
			return nil, err
		}
	}

	return ga.BuildState(opts.Owner, opts.Org, opts.Repo), nil
}
