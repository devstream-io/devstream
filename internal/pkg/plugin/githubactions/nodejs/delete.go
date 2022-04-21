package nodejs

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Delete remove GitHub Actions workflows.
func Delete(options map[string]interface{}) (bool, error) {
	var opts Options

	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return false, err
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
		return false, fmt.Errorf("opts are illegal")
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}

	log.Debugf("Language is %s.", ga.GetLanguage(opts.Language))

	for _, pipeline := range workflows {
		err := ghClient.DeleteWorkflow(pipeline, opts.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
