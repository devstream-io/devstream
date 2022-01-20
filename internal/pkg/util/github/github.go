package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var client *github.Client

func NewGithubClient() (*github.Client, error) {
	if client != nil {
		return client, nil
	}

	token := viper.GetString("github_token")
	if token == "" {
		return nil, fmt.Errorf("failed to initialize GitHub token. More info - https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}

	tc := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		),
	)

	client = github.NewClient(tc)

	return client, nil
}
