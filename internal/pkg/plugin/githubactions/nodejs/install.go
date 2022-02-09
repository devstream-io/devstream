package nodejs

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
)

// Install sets up GitHub Actions workflow(s).
func Install(options *map[string]interface{}) (bool, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return false, err
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	gitHubClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}

	log.Infof("Language is: %s.", ga.GetLanguage(opt.Language))

	for _, w := range workflows {
		if err := gitHubClient.AddWorkflow(w, opt.Branch); err != nil {
			return false, err
		}
	}

	return true, nil
}
