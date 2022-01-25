package githubactions

import (
	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/util/github"
)

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := githubActions.GetLanguage()
	log.Infof("Language is: %s.", language.String())

	// if docker is enabled, create repo secrets for DOCKERHUB_USERNAME and DOCKERHUB_TOKEN
	if githubActions.options.Jobs.Docker.Enable {
		ghOptions := &github.Option{
			Owner:    githubActions.options.Owner,
			Repo:     githubActions.options.Repo,
			NeedAuth: true,
			WorkPath: github.DefaultWorkPath,
		}
		c, err := github.NewClient(ghOptions)
		if err != nil {
			return false, err
		}
		if err := c.AddRepoSecret("DOCKERHUB_USERNAME", viper.GetString("dockerhub_username")); err != nil {
			return false, err
		}
		if err := c.AddRepoSecret("DOCKERHUB_TOKEN", viper.GetString("dockerhub_token")); err != nil {
			return false, err
		}
	}

	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())

	for _, w := range ws {
		if err := githubActions.renderTemplate(w); err != nil {
			return false, err
		}
		if err := githubActions.AddWorkflow(w); err != nil {
			return false, err
		}
	}

	return true, nil
}
