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

	// TODO(ironcore864): WorkPath should not belong to "Option",
	// because WorkPath is only used when calling the download function,
	// and it's not a property of the github client.
	WorkPath string
}

func NewClient(option *Option) (*Client, error) {
	// same option will get same client
	if client != nil && *client.Option == *option {
		log.Debug("Used a cached client")
		return client, nil
	}

	defer func() {
		if client.Option.WorkPath == "" {
			log.Debugf("Used the default workpath: %s", DefaultWorkPath)
			client.Option.WorkPath = DefaultWorkPath
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

	// TODO(ironcore864): The github package should not depend on dtm.
	// GitHub util should be a public util, instead of internal.
	// So, it should be placed under /pkg/ instead of /internal/pkg/
	// And, since this is a "util" package, it should be able to be used directly without using DTM.
	// At the moment, viper.GetString() depends on viper.BindEnv() which is triggered in the dtm main file,
	// which means, if somebody uses this package in his own package, internal or external,
	// it will fail without calling the following code first:
	//
	// if err := viper.BindEnv("github_token"); err != nil {
	// 	log.Fatal(err)
	// }
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
