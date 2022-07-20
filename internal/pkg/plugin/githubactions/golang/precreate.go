package golang

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
)

func createDockerHubInfoForPush(options plugininstaller.RawOptions) error {
	opts, err := github.NewGithubActionOptions(options)
	if err != nil {
		return err
	}
	ghClient, err := opts.GetGithubClient()
	if err != nil {
		return err
	}
	if opts.Docker != nil && opts.Docker.Enable {
		dockerhubToken := viper.GetString("dockerhub_token")
		if dockerhubToken == "" {
			return fmt.Errorf("DockerHub Token is empty")
		}

		err = ghClient.AddRepoSecret("DOCKERHUB_TOKEN", dockerhubToken)
		if err != nil {
			return err
		}
		dockerhubUserName := viper.GetString("dockerhub_username")
		err := ghClient.AddRepoSecret("DOCKERHUB_USERNAME", dockerhubUserName)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDockerHubInfoForPush(options plugininstaller.RawOptions) error {
	opts, err := github.NewGithubActionOptions(options)
	if err != nil {
		return err
	}
	ghClient, err := opts.GetGithubClient()
	if err != nil {
		return err
	}
	if opts.Docker != nil && opts.Docker.Enable {
		for _, secret := range []string{"DOCKERHUB_USERNAME", "DOCKERHUB_TOKEN"} {
			if err := ghClient.DeleteRepoSecret(secret); err != nil {
				return err
			}
		}
	}
	return nil
}
