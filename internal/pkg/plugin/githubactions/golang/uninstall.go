package golang

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
)

// Uninstall remove GitHub Actions workflows.
func Uninstall(options *map[string]interface{}) (bool, error) {
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

	// if docker is enabled, delete repo secrets DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if opt.Docker.Enable == "True" {
		for _, secret := range []string{"DOCKERHUB_USERNAME", "DOCKERHUB_TOKEN"} {
			if err := gitHubClient.DeleteRepoSecret(secret); err != nil {
				return false, err
			}
		}
	}

	for _, pipeline := range workflows {
		err := gitHubClient.DeleteWorkflow(pipeline, opt.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
