package python

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
)

// Create sets up GitHub Actions workflow(s).
func Create(options *map[string]interface{}) (map[string]interface{}, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	gitHubClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	log.Infof("Language is: %s.", ga.GetLanguage(opt.Language))

	for _, w := range workflows {
		if err := gitHubClient.AddWorkflow(w, opt.Branch); err != nil {
			return nil, err
		}
	}

	return ga.BuildState(opt.Owner, opt.Repo), nil
}
