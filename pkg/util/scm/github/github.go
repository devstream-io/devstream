package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	defaultWorkPath          = ".github-workpath"
	defaultCommitAuthor      = "devstream"
	defaultCommitAuthorEmail = "devstream@merico.dev"
	repoPlaceHolderFileName  = ".placeholder"
	TokenEnvKey              = "GITHUB_TOKEN"
)

var (
	client *Client
)

type Client struct {
	*git.RepoInfo
	*github.Client
	context.Context
}

func NewClient(option *git.RepoInfo) (*Client, error) {
	// same option will get same client
	if client != nil && *client.RepoInfo == *option {
		log.Debug("Use a cached client.")
		return client, nil
	}

	var retErr error
	defer func() {
		if retErr != nil {
			return
		}
		if client.RepoInfo.WorkPath == "" {
			client.RepoInfo.WorkPath = defaultWorkPath
		}
	}()

	// a. client without auth enabled
	if !option.NeedAuth {
		log.Debug("Auth is not enabled.")
		client = &Client{
			RepoInfo: option,
			Client:   github.NewClient(nil),
			Context:  context.Background(),
		}

		return client, nil
	}
	log.Debug("Auth is enabled.")

	// b. client with auth enabled

	// TL;DR: Don't use viper.GetString("xxx") in the `util/xxx` package.
	// Don't use `token := viper.GetString("github_token")` here,
	// it will fail without calling `viper.BindEnv("github_token")` first.
	// os.Getenv() function is more clear and reasonable here.
	token := os.Getenv(TokenEnvKey)
	if token == "" {
		// github_token works well as GITHUB_TOKEN.
		token = os.Getenv("github_token")
	}
	if token == "" {
		retErr = fmt.Errorf("environment variable GITHUB_TOKEN is not set. Failed to initialize GitHub token. More info - " +
			"https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
		return nil, retErr
	}
	log.Debugf("Token: %s.", token)

	tc := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		),
	)

	client = &Client{
		RepoInfo: option,
		Client:   github.NewClient(tc),
		Context:  context.Background(),
	}

	return client, nil
}
