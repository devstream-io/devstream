package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
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
		Org:      opts.Org,
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
		dockerhubToken := viper.GetString("dockerhub_token")
		if dockerhubToken == "" {
			return nil, fmt.Errorf("DockerHub Token is empty")
		}

		if err := ghClient.AddRepoSecret("DOCKERHUB_TOKEN", dockerhubToken); err != nil {
			return nil, err
		}
	}

	for _, w := range workflows {
		content, err := renderTemplate(w, &opts)
		if err != nil {
			return nil, err
		}
		w.WorkflowContent = content
		if err := ghClient.AddWorkflow(w, opts.Branch); err != nil {
			return nil, err
		}
	}

	return ga.BuildState(opts.Owner, opts.Org, opts.Repo), nil
}
