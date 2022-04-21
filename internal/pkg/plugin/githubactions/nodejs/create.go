package nodejs

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create sets up GitHub Actions workflow(s).
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	log.Debugf("Language is: %s.", ga.GetLanguage(opts.Language))

	for _, w := range workflows {
		if err := ghClient.AddWorkflow(w, opts.Branch); err != nil {
			return nil, err
		}
	}

	return ga.BuildState(opts.Owner, opts.Org, opts.Repo), nil
}
