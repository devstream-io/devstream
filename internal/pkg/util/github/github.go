package github

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/merico-dev/stream/internal/pkg/log"
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

func (c *Client) DownloadAsset(tagName, assetName string) error {
	// 1. get releases
	releases, resp, err := c.Repositories.ListReleases(context.TODO(), c.Owner, c.Repo, &github.ListOptions{})
	if err != nil {
		return err
	}
	log.Debug("Got releases successful.")
	for i, r := range releases {
		log.Debugf("Release(%d): %s", i+1, r.GetName())
	}

	log.Debugf("Response status: %s", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response status not expected: %s", resp.Status)
	}

	// 2. get assets
	var assets []*github.ReleaseAsset
	for _, r := range releases {
		if *r.TagName != tagName {
			continue
		}
		log.Debugf("Got a matched tag %s with release <%s>", *r.TagName, *r.Name)

		if len(r.Assets) == 0 {
			log.Debug("Assets is empty")
			return fmt.Errorf("assets is empty")
		}
		log.Debugf("%d Assets was found", len(r.Assets))

		assets = r.Assets
		break
	}
	if len(assets) == 0 {
		log.Debugf("Release with tag <%s> was not found", tagName)
		return fmt.Errorf("release with tag <%s> was not found", tagName)
	}

	// 3. get download url
	// format: https://github.com/merico-dev/dtm-scaffolding-golang/releases/download/v0.0.1/dtm-scaffolding-golang-v0.0.1.tar.gz
	var downloadUrl string
	for _, a := range assets {
		if a.GetName() == assetName {
			downloadUrl = a.GetBrowserDownloadURL()
			log.Debugf("Download url: %s", downloadUrl)
			break
		}
	}
	if downloadUrl == "" {
		log.Debugf("Failed to got the download url for %s, maybe it not exists", assetName)
		return fmt.Errorf("failed to got the download url for %s, maybe it not exists", assetName)
	}

	// 4. download
	n, err := downloader.Download(downloadUrl, c.WorkPath)
	if err != nil {
		log.Debugf("Failed to download asset from %s", downloadUrl)
		return err
	}
	log.Debugf("Downloaded <%d> bytes", n)

	return nil
}

func (c *Client) InitRepoLocalAndPushToRemote(repoPath string) error {
	err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			log.Debugf("Found dir: %s", path)
			return nil
		}

		log.Debugf("Found file: %s", path)
		return c.CreateFile(path)
	})

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateFile(filePath string) error {
	_, _, err := c.Repositories.CreateFile(context.TODO(), c.Owner, c.Repo, filePath, &github.RepositoryContentFileOptions{})
	return err
}
