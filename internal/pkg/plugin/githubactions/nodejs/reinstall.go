package nodejs

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
)

// Reinstall remove and set up GitHub Actions workflows.
func Reinstall(options *map[string]interface{}) (bool, error) {
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

	log.Infof("language is %s", ga.GetLanguage(opt.Language))

	for _, pipeline := range workflows {
		err := gitHubClient.DeleteWorkflow(pipeline, opt.Branch)
		if err != nil {
			return false, err
		}

		err = gitHubClient.AddWorkflow(pipeline, opt.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
