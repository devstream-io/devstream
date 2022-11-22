package python

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/github"
)

func createDockerHubInfoForPush(options configmanager.RawOptions) error {
	opts, err := github.NewGithubActionOptions(options)
	if err != nil {
		return err
	}
	ghClient, err := opts.GetGithubClient()
	if err != nil {
		return err
	}

	if opts.CheckAddDockerHubToken() {
		dockerhubToken := viper.GetString("dockerhub_token")
		if dockerhubToken == "" {
			return fmt.Errorf("DockerHub Token is empty")
		}

		err = ghClient.AddRepoSecret("DOCKERHUB_TOKEN", dockerhubToken)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDockerHubInfoForPush(options configmanager.RawOptions) error {
	opts, err := github.NewGithubActionOptions(options)
	if err != nil {
		return err
	}
	ghClient, err := opts.GetGithubClient()
	if err != nil {
		return err
	}
	if opts.Docker != nil && opts.Docker.Enable {
		for _, secret := range []string{"DOCKERHUB_TOKEN"} {
			if err := ghClient.DeleteRepoSecret(secret); err != nil {
				return err
			}
		}
	}
	return nil
}
