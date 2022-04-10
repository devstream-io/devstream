package golang

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

	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
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

	log.Debugf("language is %s.", ga.GetLanguage(opts.Language))

	// if docker is enabled, delete repo secrets DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if opts.Docker != nil && opts.Docker.Enable {
		for _, secret := range []string{"DOCKERHUB_USERNAME", "DOCKERHUB_TOKEN"} {
			if err := ghClient.DeleteRepoSecret(secret); err != nil {
				return false, err
			}
		}
	}

	for _, pipeline := range workflows {
		err := ghClient.DeleteWorkflow(pipeline, opts.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
