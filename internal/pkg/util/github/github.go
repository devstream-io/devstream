package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/merico-dev/stream/internal/pkg/util/downloader"
)

const DefaultWorkPath = ".github"

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
		return client, nil
	}

	defer func() {
		if client.WorkPath == "" {
			client.WorkPath = DefaultWorkPath
		}
	}()

	// client without auth enabled
	if !option.NeedAuth {
		client = &Client{
			Option: option,
			Client: github.NewClient(nil),
		}

		return client, nil
	}

	// client with auth enabled
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

	client = &Client{
		Option: option,
		Client: github.NewClient(tc),
	}

	return client, nil
}

func (c *Client) DownloadAsset(releaseName, assetName string) error {
	// 1. get releases
	releases, resp, err := c.Repositories.ListReleases(context.TODO(), c.Owner, c.Repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response status not expected: %s", resp.Status)
	}

	// 2. get assets
	var assets []*github.ReleaseAsset
	for _, r := range releases {
		if *r.Name != releaseName {
			continue
		}

		if len(r.Assets) == 0 {
			return fmt.Errorf("assets is empty")
		}

		assets = r.Assets
		break
	}
	if len(assets) == 0 {
		return fmt.Errorf("r not found: %s", releaseName)
	}

	// 3. get download url
	// format: https://github.com/merico-dev/stream/releases/download/v0.0.1/argocdapp_0.0.1.so
	var downloadUrl string
	for _, a := range assets {
		if a.GetName() == assetName {
			downloadUrl = a.GetBrowserDownloadURL()
		}
	}
	if downloadUrl == "" {
		return fmt.Errorf("failed to got the download url for %s, maybe it not exists", assetName)
	}

	// 4. download
	_, err = downloader.Download(downloadUrl, c.WorkPath)
	if err != nil {
		return err
	}

	return nil
}
