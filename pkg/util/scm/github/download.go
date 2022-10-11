package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) DownloadAsset(tagName, assetName, fileName string) error {
	// 1. get releases
	releases, _, err := c.Repositories.ListReleases(context.TODO(), c.GetRepoOwner(), c.Repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	log.Debug("Got releases successful.")
	for i, r := range releases {
		log.Debugf("Release(%d): %s.", i+1, r.GetName())
	}

	// 2. get assets
	var assets []*github.ReleaseAsset
	for _, r := range releases {
		if *r.TagName != tagName {
			continue
		}
		log.Debugf("Got a matched tag %s with release <%s>.", *r.TagName, *r.Name)

		if len(r.Assets) == 0 {
			log.Debug("Assets is empty.")
			return fmt.Errorf("assets is empty")
		}
		log.Debugf("%d Assets was found.", len(r.Assets))

		assets = r.Assets
		break
	}
	if len(assets) == 0 {
		log.Debugf("Release with tag <%s> was not found.", tagName)
		return fmt.Errorf("release with tag <%s> was not found", tagName)
	}

	// 3. get download url
	// format: https://github.com/merico-dev/dtm-scaffolding-golang/releases/download/v0.0.1/dtm-scaffolding-golang-v0.0.1.tar.gz
	var downloadUrl string
	for _, a := range assets {
		if a.GetName() == assetName {
			downloadUrl = a.GetBrowserDownloadURL()
			log.Debugf("Download url: %s.", downloadUrl)
			break
		}
	}
	if downloadUrl == "" {
		log.Debugf("Failed to got the download url for %s, maybe it not exists.", assetName)
		return fmt.Errorf("failed to got the download url for %s, maybe it not exists", assetName)
	}

	// 4. download
	n, err := downloader.New().WithProgressBar().Download(downloadUrl, fileName, c.WorkPath)
	if err != nil {
		log.Debugf("Failed to download asset from %s.", downloadUrl)
		return err
	}
	log.Debugf("Downloaded <%d> bytes.", n)

	return nil
}
