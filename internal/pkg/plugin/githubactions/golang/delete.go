package golang

import (
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Delete remove GitHub Actions workflows.
func Delete(options map[string]interface{}) (bool, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return false, err
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return false, err
	}

	log.Debugf("language is %s.", ga.GetLanguage(opt.Language))

	// if docker is enabled, delete repo secrets DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if opt.Docker != nil && opt.Docker.Enable {
		for _, secret := range []string{"DOCKERHUB_USERNAME", "DOCKERHUB_TOKEN"} {
			if err := ghClient.DeleteRepoSecret(secret); err != nil {
				return false, err
			}
		}
	}

	for _, pipeline := range workflows {
		err := ghClient.DeleteWorkflow(pipeline, opt.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
