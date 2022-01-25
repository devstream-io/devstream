package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/merico-dev/stream/internal/pkg/log"
)

const (
	DefaultWorkPath = ".github"
	// https://github.com/merico-dev/dtm-scaffolding-golang/archive/refs/heads/main.zip -> 302 ->
	// https://codeload.github.com/merico-dev/dtm-scaffolding-golang/zip/refs/heads/main
	DefaultLatestCodeZipfileDownloadUrlFormat = "https://codeload.github.com/%s/%s/zip/refs/heads/main"
	DefaultLatestCodeZipfileName              = "main-latest.zip"
)

var client *Client

type Client struct {
	*Option
	*github.Client
}

type Option struct {
	Owner    string
	Repo     string
	NeedAuth bool
	// default -> ".github"
	WorkPath string
}

func NewClient(option *Option) (*Client, error) {
	// same option will get same client
	if client != nil && *client.Option == *option {
		log.Debug("Used a cached client")
		return client, nil
	}

	defer func() {
		if client.WorkPath == "" {
			log.Debugf("Used the default workpath: %s", DefaultWorkPath)
			client.WorkPath = DefaultWorkPath
		}
	}()

	// client without auth enabled
	if !option.NeedAuth {
		log.Debug("Auth is not enabled")
		client = &Client{
			Option: option,
			Client: github.NewClient(nil),
		}

		return client, nil
	}
	log.Debug("Auth is enabled")

	// client with auth enabled
	token := viper.GetString("github_token")
	if token == "" {
		return nil, fmt.Errorf("failed to initialize GitHub token. More info - https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}
	log.Debugf("Token: %s", token)

	tc := oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		),
	)

	client = &Client{
		Option: option,
		Client: github.NewClient(tc),
	}

	return client, nil
}
