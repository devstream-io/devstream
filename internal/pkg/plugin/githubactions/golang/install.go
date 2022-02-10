package golang

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/spf13/viper"
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

	// if docker is enabled, create repo secrets for DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if opt.Docker.Enable {
		if err := gitHubClient.AddRepoSecret("DOCKERHUB_USERNAME", viper.GetString("dockerhub_username")); err != nil {
			return false, err
		}
		if err := gitHubClient.AddRepoSecret("DOCKERHUB_TOKEN", viper.GetString("dockerhub_token")); err != nil {
			return false, err
		}
	}

	for _, w := range workflows {
		content, err := renderTemplate(w, opt)
		if err != nil {
			return false, err
		}
		w.WorkflowContent = content
		if err := gitHubClient.AddWorkflow(w, opt.Branch); err != nil {
			return false, err
		}
	}

	return true, nil
}
