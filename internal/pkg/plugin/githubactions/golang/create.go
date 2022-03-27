package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create sets up GitHub Actions workflow(s).
func Create(options map[string]interface{}) (map[string]interface{}, error) {

	var opts Options

	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
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

	log.Debugf("Language is: %s.", ga.GetLanguage(opts.Language))

	// if docker is enabled, create repo secrets for DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if opts.Docker != nil && opts.Docker.Enable {
		if err := ghClient.AddRepoSecret("DOCKERHUB_USERNAME", viper.GetString("dockerhub_username")); err != nil {
			return nil, err
		}
		if err := ghClient.AddRepoSecret("DOCKERHUB_TOKEN", viper.GetString("dockerhub_token")); err != nil {
			return nil, err
		}
	}

	for _, w := range workflows {
		content, err := renderTemplate(w, opts)
		if err != nil {
			return nil, err
		}
		w.WorkflowContent = content
		if err := ghClient.AddWorkflow(w, opts.Branch); err != nil {
			return nil, err
		}
	}

	return ga.BuildState(opts.Owner, opts.Repo), nil
}
