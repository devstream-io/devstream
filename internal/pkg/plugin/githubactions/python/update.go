package python

import (
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Update remove and set up GitHub Actions workflows.
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

	return ga.BuildState(opts.Owner, opts.Repo), nil
}
