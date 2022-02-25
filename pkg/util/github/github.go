package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"

	"github.com/merico-dev/stream/pkg/util/log"
)

const (
	DefaultWorkPath = ".github-workpath"
	// https://github.com/merico-dev/dtm-scaffolding-golang/archive/refs/heads/main.zip -> 302 ->
	// https://codeload.github.com/merico-dev/dtm-scaffolding-golang/zip/refs/heads/main
	DefaultLatestCodeZipfileDownloadUrlFormat = "https://codeload.github.com/%s/%s/zip/refs/heads/main"
	DefaultLatestCodeZipfileName              = "main-latest.zip"
)

var client *Client

type Client struct {
	*Option
	*github.Client
	context.Context
}

type Option struct {
	Owner    string
	Repo     string
	NeedAuth bool
	WorkPath string
}

func NewClient(option *Option) (*Client, error) {
	// same option will get same client
	if client != nil && *client.Option == *option {
		log.Debug("Use a cached client.")
		return client, nil
	}

	defer func() {
		if client.Option.WorkPath == "" {
			client.Option.WorkPath = DefaultWorkPath
		}
	}()

	// a. client without auth enabled
	if !option.NeedAuth {
		log.Debug("Auth is not enabled.")
		client = &Client{
			Option:  option,
			Client:  github.NewClient(nil),
			Context: context.Background(),
		}

		return client, nil
	}
	log.Debug("Auth is enabled.")

	// b. client with auth enabled

	// Don't use `token := viper.GetString("github_token")` here,
	// it will fail without calling `viper.BindEnv("github_token")` first.
	// os.Getenv() function is more clear and reasonable here.
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// github_token works well as GITHUB_TOKEN.
		token = os.Getenv("github_token")
	}
	if token == "" {
		return nil, fmt.Errorf("failed to initialize GitHub token. More info - " +
			"https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}
	log.Debugf("Token: %s.", token)

	ctx := context.Background()
	tc := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		),
	)

	client = &Client{
		Option:  option,
		Client:  github.NewClient(tc),
		Context: ctx,
	}

	return client, nil
}
